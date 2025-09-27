-- Create the stock prices table if it doesn't exist
CREATE TABLE IF NOT EXISTS stock_prices (
  symbol TEXT NOT NULL,
  open NUMERIC(12, 4),
  close NUMERIC(12, 4),
  high NUMERIC(12, 4),
  low NUMERIC(12, 4),
  timestamp TIMESTAMPTZ NOT NULL,
  volume BIGINT,
  PRIMARY KEY (symbol, timestamp)
);

-- Convert to hypertable if not already (separate transaction)
SELECT create_hypertable('stock_prices', 'timestamp', if_not_exists => TRUE);
