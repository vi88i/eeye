# eeye
eagle eye stock screener

# Installation

## Prerequisites

1. Go 1.24.5 or later
2. Docker (for running the database)

## Installing Dependencies

1. Clone the repository:
```bash
git clone https://github.com/vi88i/eeye.git
cd eeye
```

2. Install Go dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
# Copy the environment template
cp .env.example .env

# Edit .env and update only the GROWW_ACCESS_TOKEN
# You'll need to set:
# - GROWW_ACCESS_TOKEN: Your Groww API access token
# All other variables are pre-configured with their default values
```

# Development Setup

## Database Setup

```cmd
docker pull timescale/timescaledb:latest-pg14

docker run -d --name eeye-db `
  -p 5432:5432 `
  -v eeye-vol:/var/lib/postgresql/data `
  -e POSTGRES_PASSWORD=root `
  -e POSTGRES_USER=admin `
  -e POSTGRES_DB=eeye `
  timescale/timescaledb:latest-pg14

docker exec -it eeye-db psql -U admin -d eeye
```

## Pre-commit Hooks

This project uses pre-commit hooks to ensure code quality. The hooks will:
- Format Go code using `go fmt`
- Run `golangci-lint` for code linting

### Prerequisites

1. Install golangci-lint:
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Installing the Pre-commit Hooks

To install the pre-commit hooks:

1. Make sure you're in the root directory of the project
2. Run the installation script:
```bash
# On Unix/Linux/MacOS
chmod +x scripts/install-hooks.sh
./scripts/install-hooks.sh

# On Windows (Git Bash or similar)
sh scripts/install-hooks.sh
```

The hooks will now run automatically before each commit. If there are any formatting issues or linting errors, the commit will be blocked until they are fixed.
