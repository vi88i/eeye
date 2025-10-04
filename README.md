# eeye
`e`agle `eye` stock screener

# Overview

`eeye` is a stock screener that uses various technical analysis strategies to identify potential trading opportunities in the stock market. It leverages the Groww API to fetch real-time market data and applies a series of screening techniques to filter stocks based on user-defined criteria.

## Features

- Fetching all available stocks from NSE
- Integration with Groww API for real-time stock data
- Multiple technical analysis strategies (e.g., Bollinger Bands, EMA, RSI)
- Modular design for easy addition of new strategies

## How it works?

1. **Fetch Stock Data**: The application fetches historical stock data from the Groww API.
2. **Apply Screening Steps**: Each stock is processed through a series of screening steps defined in the strategy executor. These steps include various technical analysis techniques.
3. **Output Results**: Stocks that pass all screening steps are printed to the console as potential trading opportunities.

# Installation

## Prerequisites

1. Go 1.24.5 or later
2. Docker (for running the database)

## Installing Dependencies

1. Run the setup script:
```bash
# On Unix/Linux/MacOS
chmod +x scripts/setup-env.sh
./scripts/setup-env.sh

# On Windows (Git Bash or similar)
sh scripts/setup-env.sh
```

2. Configure your Groww API token:
```bash
# Edit .env and set your GROWW_ACCESS_TOKEN
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

# Running the Application

To start the application, use the provided start script:

```bash
# On Unix/Linux/MacOS
chmod +x scripts/start.sh
./scripts/start.sh

# On Windows (Git Bash or similar)
sh scripts/start.sh
```

The start script will:
1. Check if the database container is running and start it if needed
2. Determine the correct project paths
3. Run the application with `go run`

Note: Make sure you've completed the database setup and configuration steps before running the application.
