#!/bin/bash

# Exit on error
set -e

# Detect OS and Git Bash
case "$(uname -s)" in
    Linux*)     OS="Linux";;
    Darwin*)    OS="Mac";;
    MINGW*)     OS="GitBash";;
    CYGWIN*)    OS="Cygwin";;
    *)          OS="Unknown";;
esac

# Function to normalize path for Docker
normalize_path() {
    local path=$1
    
    # If running in Git Bash or Cygwin on Windows
    if [ "$OS" = "GitBash" ] || [ "$OS" = "Cygwin" ]; then
        # Convert Windows path to Unix style
        path="$(cygpath -u "$path")"
    fi
    
    echo "$path"
}

# Get absolute path of script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "Operating System: $OS"
echo "Project Root: $PROJECT_ROOT"

# Check if Docker is available
if ! command -v docker >/dev/null 2>&1; then
    echo "Error: Docker is not installed or not in PATH"
    exit 1
fi

# Function to setup Docker container
setup_container() {
    echo "Setting up TimescaleDB container..."
    
    # Check if image exists locally
    if ! docker image inspect timescale/timescaledb:latest-pg14 >/dev/null 2>&1; then
        echo "Pulling TimescaleDB image..."
        docker pull timescale/timescaledb:latest-pg14
    else
        echo "TimescaleDB image already exists locally"
    fi

    # Check if volume exists, create if it doesn't
    if ! docker volume inspect eeye-vol >/dev/null 2>&1; then
        echo "Creating Docker volume eeye-vol..."
        docker volume create eeye-vol
    else
        echo "Docker volume eeye-vol already exists"
    fi

    # Check if container exists
    if docker ps -a | grep -q eeye-db; then
        if docker ps | grep -q eeye-db; then
            echo "Container eeye-db is already running"
        else
            echo "Starting existing container eeye-db..."
            docker start eeye-db
        fi
    else
        echo "Creating and starting new container eeye-db..."
        docker run -d --name eeye-db \
            -p 5432:5432 \
            -v eeye-vol:/var/lib/postgresql/data \
            -e POSTGRES_PASSWORD=root \
            -e POSTGRES_USER=admin \
            -e POSTGRES_DB=eeye \
            timescale/timescaledb:latest-pg14
    fi

    # Verify the container is running
    if ! docker ps | grep -q eeye-db; then
        echo "Error: Failed to start container"
        exit 1
    fi
    
    echo "TimescaleDB container is ready"
}

# Setup the container
setup_container

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
until docker exec eeye-db pg_isready -U admin > /dev/null 2>&1; do
    echo "PostgreSQL is unavailable - sleeping for 1 second"
    sleep 1
done

# Check if TimescaleDB extension is available and create it
echo "Setting up TimescaleDB extension..."
docker exec eeye-db psql -U admin -d eeye -c "CREATE EXTENSION IF NOT EXISTS timescaledb;" || {
    echo "Error: Failed to create TimescaleDB extension"
    exit 1
}

# Apply schema
echo "Applying database schema..."

# Construct and normalize schema path
SCHEMA_PATH="$PROJECT_ROOT/sql/SCHEMA.sql"
NORMALIZED_SCHEMA_PATH=$(normalize_path "$SCHEMA_PATH")

echo "Looking for schema at: $NORMALIZED_SCHEMA_PATH"

if [ ! -f "$NORMALIZED_SCHEMA_PATH" ]; then
    echo "Error: SCHEMA.sql not found at $NORMALIZED_SCHEMA_PATH"
    exit 1
fi

echo "Applying schema to database..."
# Read the schema file and pipe it directly to psql in the container
if ! cat "$NORMALIZED_SCHEMA_PATH" | docker exec -i eeye-db psql -U admin -d eeye; then
    echo "Error: Failed to apply schema"
    exit 1
fi

echo "Database setup completed successfully!"
