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

The eeye stock screener operates through a multi-stage pipeline that efficiently processes thousands of stocks in parallel:

### 1. Data Acquisition Phase

**NSE Stock Discovery**
- Downloads the latest NSE Bhavcopy (daily market report) containing all listed stocks
- Filters stocks to include only equity shares (EQ series) that are currently listed
- Identifies the last trading day to determine data freshness

**Database Synchronization**
- Compares fetched stocks with existing database records
- Identifies three categories of stocks:
  - **Newly listed stocks**: Never seen before in the database
  - **Out-of-sync stocks**: Missing recent trading data
  - **De-listed stocks**: No longer present in NSE data (cleaned up with `--cleanup` flag)

**Historical Data Backfill**
- Uses multiple worker goroutines to fetch missing historical data in parallel
- Respects API rate limits (configurable requests per second)
- Fetches OHLCV (Open, High, Low, Close, Volume) data from Groww API
- Stores data in TimescaleDB hypertable for efficient time-series queries

### 2. Analysis Phase

**In-Memory Caching**
- Loads required stock data into memory cache for fast access
- Avoids repeated database queries during analysis

**Parallel Strategy Execution**
- Spawns multiple worker goroutines to process stocks concurrently
- Each stock is evaluated against all configured strategies simultaneously
- Strategies are completely independent and run in parallel

**Strategy Screening Pipeline**

Each strategy consists of multiple screening steps that must ALL pass. Stocks are evaluated using technical indicators (EMA, RSI, Bollinger Bands, Volume MA), pattern recognition, and custom logic specific to each strategy.

### 3. Results Aggregation

**Collection**
- Each strategy collects stocks that pass all its screening steps
- Results are aggregated from all parallel workers

**Output**
- Prints matching stock symbols grouped by strategy
- Logs execution time and performance metrics

### 4. Optional Modes

**MCP Server Mode** (`--mcp` flag)
- Runs as an HTTP server providing Model Context Protocol interface
- Exposes tools for querying technical data and OHLC information
- Allows AI assistants like Claude to analyze stock data interactively
- No automated screening in this mode - responds to queries on demand

**Cleanup Mode** (`--cleanup` flag)
- Removes data for de-listed stocks after analysis completes
- Keeps database size manageable and data relevant

### Architecture Highlights

```
┌─────────────────────────────────────────────────────────────┐
│                     NSE Bhavcopy API                        │
└────────────────────────┬────────────────────────────────────┘
                         │ Download
                         ▼
              ┌──────────────────────┐
              │   Stock Discovery    │
              └──────────┬───────────┘
                         │ Identify gaps
                         ▼
              ┌──────────────────────┐
              │  Parallel Ingestion  │◄──── Groww API (OHLCV data)
              │   (Worker Pool)      │
              └──────────┬───────────┘
                         │ Store
                         ▼
              ┌──────────────────────┐
              │   TimescaleDB        │
              │   (Hypertable)       │
              └──────────┬───────────┘
                         │ Load to cache
                         ▼
              ┌──────────────────────┐
              │   In-Memory Cache    │
              └──────────┬───────────┘
                         │ Feed stocks
                         ▼
        ┌────────────────────────────────────┐
        │   Strategy Workers (Parallel)      │
        │  ┌──────────────────────────────┐  │
        │  │  Strategy 1: Bullish Swing   │  │
        │  │  Strategy 2: EMA Breakdown   │  │
        │  │  Strategy 3: RSI Momentum    │  │
        │  │  Strategy 4: BB Reversal     │  │
        │  │  Strategy 5: Breakout        │  │
        │  └──────────────────────────────┘  │
        └────────────────┬───────────────────┘
                         │ Collect results
                         ▼
              ┌──────────────────────┐
              │  Results Aggregator  │
              └──────────┬───────────┘
                         │ Print
                         ▼
                    Console Output
```

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

## Running as Stock Screener

To start the application as a stock screener, use the provided start script:

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

## Command Line Flags

The application supports the following command-line flags:

- `--mcp`: Enable MCP (Model Context Protocol) server mode
- `--cleanup`: Clean up de-listed stocks from the database after analysis

### Examples

```bash
# Run stock screener and clean up de-listed stocks
go run main.go --cleanup

# Start MCP server
go run main.go --mcp

# Run screener without cleanup (default)
go run main.go
```

## Running as MCP Server

The application can run as an MCP (Model Context Protocol) server, allowing AI assistants like Claude to analyze stock data through a standardized interface.

### Starting the MCP Server

```bash
go run main.go --mcp
```

The server will start on the host and port specified in your `.env` file (defaults: `localhost:3000`).

### Configuring Claude Desktop

To use the eeye MCP server with Claude Desktop, you need to:

1. **Start the MCP server** (in a separate terminal):
   ```bash
   go run main.go --mcp
   ```

   The server will start on `http://localhost:3000` (or the host/port configured in your `.env` file).

2. **Configure Claude Desktop** to connect to the HTTP server:

**Location of config file:**
- macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
- Windows: `%APPDATA%\Claude\claude_desktop_config.json`

**Configuration:**

```json
{
  "mcpServers": {
    "eeye": {
      "command": "npx",
      "args": [
        "mcp-remote",
        "http://localhost:3000/"
      ]
    }
  }
}
```

**Note:** Make sure the port in the URL matches your `MCP_PORT` setting in the `.env` file (default is 3000).

**After updating the configuration:**
1. Restart Claude Desktop
2. The eeye MCP server should be running in a separate terminal
3. The eeye tools and resources will be available in Claude Desktop

**Troubleshooting:**
- Make sure the eeye server is running before starting Claude Desktop
- Check that the port number matches between your `.env` file and Claude Desktop config
- If you get connection errors, verify the server is accessible at `http://localhost:3000/` in your browser

### Available MCP Resources

- **nseStocks** (`db:stocks`): Returns a comma-separated list of all NSE stock symbols available in the database

### Available MCP Tools

1. **getTechnicalData**
   - **Description**: Provides comprehensive technical analysis data including OHLC, EMA (5, 13, 26, 50), RSI, and volume indicators
   - **Input**: `{ "symbol": "STOCK_SYMBOL" }`
   - **Output**: Array of technical data sorted by date (most recent first)

2. **getOhlcData**
   - **Description**: Provides basic OHLC (Open, High, Low, Close) data with timestamps
   - **Input**: `{ "symbol": "STOCK_SYMBOL" }`
   - **Output**: Array of OHLC data sorted by date (most recent first)

### Example Prompts for Claude

Once configured, you can ask Claude questions like:

1. **Get available stocks:**
   ```
   What NSE stocks are available in the eeye database?
   ```

2. **Technical analysis:**
   ```
   Get the technical data for ZOMATO and analyze the recent trend.
   Is ZOMATO showing bullish momentum based on its EMA and RSI indicators?
   ```

3. **Price analysis:**
   ```
   Get the OHLC data for RELIANCE for the last 30 days and identify support/resistance levels.
   ```

4. **Compare stocks:**
   ```
   Compare the RSI and EMA50 indicators between INFY and TCS. Which one looks better positioned?
   ```

5. **Pattern recognition:**
   ```
   Analyze the OHLC data for TATAMOTORS and identify any bullish or bearish patterns in the last 10 trading days.
   ```

6. **Multi-stock screening:**
   ```
   Get the list of available stocks, then check which ones have RSI between 40-60 and are trading above their EMA50.
   ```

---

⚠️ Disclaimer

This app is for educational purposes only and does not provide financial advice. Use at your own risk. Always do your own research before making investment decisions.
