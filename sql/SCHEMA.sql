CREATE EXTENSION IF NOT EXISTS timescaledb;

CREATE TABLE stock_prices (
  symbol TEXT NOT NULL,
  open NUMERIC(12, 4),
  close NUMERIC(12, 4),
  high NUMERIC(12, 4),
  low NUMERIC(12, 4),
  timestamp TIMESTAMPTZ NOT NULL,
  volume BIGINT,
  PRIMARY KEY (symbol, timestamp)
);

SELECT create_hypertable('stock_prices', 'timestamp');
