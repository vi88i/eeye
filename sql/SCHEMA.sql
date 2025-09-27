-- Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb;

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

-- Convert to hypertable if not already
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM timescaledb_information.hypertables
    WHERE hypertable_name = 'stock_prices'
  ) THEN
    PERFORM create_hypertable('stock_prices', 'timestamp');
  END IF;
END $$;
