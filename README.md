# eeye
eagle eye stock screener

# Overview

eeye is a stock screener that uses various technical analysis strategies to identify potential trading opportunities in the stock market. It leverages the Groww API to fetch real-time market data and applies a series of screening techniques to filter stocks based on user-defined criteria.

## Features

- Integration with Groww API for real-time stock data
- Multiple technical analysis strategies (e.g., Bollinger Bands, EMA, RSI)
- Modular design for easy addition of new strategies

## Workflow

- Add the stock to be screened in `examples/stocks.yml`
- Run the application using `go run src/main.go`
- View the results in the console output

## How it works?

1. **Fetch Stock Data**: The application fetches historical stock data from the Groww API.
2. **Apply Screening Steps**: Each stock is processed through a series of screening steps defined in the strategy executor. These steps include various technical analysis techniques.
3. **Output Results**: Stocks that pass all screening steps are printed to the console as potential trading opportunities.

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

Run the database setup script to create and configure the TimescaleDB container:

```bash
# On Unix/Linux/MacOS
chmod +x scripts/setup-db.sh
./scripts/setup-db.sh

# On Windows (Git Bash or similar)
sh scripts/setup-db.sh
```

The setup script will:
- Check if the Docker container is running
- Wait for the database to be ready
- Verify TimescaleDB extension is available
- Create the required tables and hypertables

You can run this script multiple times safely - it will not duplicate or overwrite existing data.

Note: If you see an error about the container not running, make sure you've completed step 1 successfully.

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
