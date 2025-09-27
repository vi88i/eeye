#!/bin/sh

# Check if .git directory exists
if [ ! -d ".git" ]; then
    echo "Error: .git directory not found. Please run this script from the root of the repository."
    exit 1
fi

# Create hooks directory if it doesn't exist
mkdir -p .git/hooks

# Check if golangci-lint is installed, install if not present
if ! command -v golangci-lint &> /dev/null; then
    echo "Installing golangci-lint..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    
    # Verify installation
    if ! command -v golangci-lint &> /dev/null; then
        echo "Error: Failed to install golangci-lint. Please install it manually:"
        echo "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
        exit 1
    fi
    echo "golangci-lint installed successfully!"
fi

# Copy pre-commit hook and make it executable
cp scripts/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit

echo "Pre-commit hook installed successfully!"