package db

import (
	"context"
	"eeye/src/config"
	"eeye/src/constants"
	"eeye/src/models"
	"eeye/src/utils"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
)

// GetLastCandle retrieves the most recent candlestick data for a given stock symbol.
// The timestamp in the returned candle is adjusted to the timezone specified in DB.
func GetLastCandle(symbol string) (models.Candle, error) {
	log.Printf("getting last candle for %s\n", symbol)
	ctx := context.Background()

	// Trading API works in current timezone so do the conversion of timestamp
	rows, err := Pool.Query(ctx, `
		SELECT (timestamp AT TIME ZONE $2) as timestamp
		FROM stock_prices
		WHERE symbol = $1
		ORDER BY timestamp DESC
		LIMIT 1
	`, symbol, config.DB.Tz)

	var ret = models.Candle{
		Symbol:    symbol,
		Timestamp: time.Now().UTC().Truncate(24*time.Hour).AddDate(0, 0, -constants.LookBackDays),
	}

	if err != nil {
		return ret, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&ret.Timestamp)
		if err != nil {
			return ret, fmt.Errorf("scanning failed: %w", err)
		}
	}

	return ret, nil
}

// BackfillCandles efficiently inserts multiple candlestick records into the database
// using PostgreSQL's COPY protocol. This is optimized for bulk insertions of historical data.
func BackfillCandles(stock *models.Stock, candles []models.Candle) error {
	log.Printf("backfilling %d candles for %v\n", len(candles), stock.Symbol)
	entries := make([][]any, 0, len(candles))
	for i := range candles {
		candle := &candles[i]
		entries = append(entries, []any{
			candle.Symbol,
			candle.Open,
			candle.Close,
			candle.High,
			candle.Low,
			candle.Timestamp,
			candle.Volume,
		})
	}

	var (
		columns   = []string{"symbol", "open", "close", "high", "low", "timestamp", "volume"}
		tableName = "stock_prices"
		ctx       = context.Background()
	)

	/*
		COPY FROM is a PostgreSQL protocol (binary) which helps in efficient insertion.
		Instead of creating and closing HTTP connection per insert, it creates a single connection,
		and insert params are streamed in batches.
	*/
	_, err := Pool.CopyFrom(ctx, pgx.Identifier{tableName}, columns, pgx.CopyFromRows(entries))
	if err != nil {
		return fmt.Errorf("copy from failed: %w", err)
	}

	return nil
}

// FetchAllCandles retrieves all stored candlestick data for a given stock.
// The timestamps in the returned candles are adjusted to the timezone specified in DB.
func FetchAllCandles(stock *models.Stock) ([]models.Candle, error) {
	log.Printf("fetching all candles: %v\n", stock.Symbol)
	ctx := context.Background()

	rows, err := Pool.Query(ctx, `
		SELECT symbol, open, close, high, low, (timestamp AT TIME ZONE $2) as timestamp, volume
		FROM stock_prices
		WHERE symbol = $1
		ORDER BY timestamp ASC
	`, stock.Symbol, config.DB.Tz)

	var (
		empty = utils.EmptySlice[models.Candle]()
		res   = make([]models.Candle, 0, constants.LookBackDays)
	)

	if err != nil {
		return empty, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		candle := models.Candle{}

		err := rows.Scan(
			&candle.Symbol,
			&candle.Open,
			&candle.Close,
			&candle.High,
			&candle.Low,
			&candle.Timestamp,
			&candle.Volume,
		)

		if err != nil {
			return empty, fmt.Errorf("scanning failed: %w", err)
		}

		res = append(res, candle)
	}

	return res, nil
}

// FetchAllStocks returns distinct stocks from the DB
func FetchAllStocks() ([]models.Stock, error) {
	log.Println("fetching all distinct stocks from DB")
	ctx := context.Background()

	rows, err := Pool.Query(ctx, `
		SELECT symbol
		FROM stock_prices
		GROUP BY symbol
		ORDER BY symbol ASC
	`)

	var (
		empty = utils.EmptySlice[models.Stock]()
		res   = make([]models.Stock, 0, constants.NumOfStocks)
	)

	if err != nil {
		return empty, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		stock := models.Stock{Segment: "CASH", Exchange: "NSE"}

		err := rows.Scan(&stock.Symbol)

		if err != nil {
			return empty, fmt.Errorf("scanning failed: %w", err)
		}
		stock.Name = stock.Symbol

		res = append(res, stock)
	}

	log.Printf("fetched %v distinct stock from DB\n", len(res))
	return res, nil
}

// FetchOutOfSyncStock fetches stocks that are not synced with latest market data
func FetchOutOfSyncStock(lastTradingDay string) ([]models.Stock, error) {
	log.Println("fetching out of sync stocks")
	ctx := context.Background()

	query := fmt.Sprintf(`
		SELECT symbol
		FROM stock_prices
		GROUP BY symbol
		HAVING MAX(timestamp) <> TIMESTAMPTZ '%v'
		ORDER BY symbol ASC
	`, fmt.Sprintf("%v 00:00:00 %v", lastTradingDay, config.DB.Tz))
	rows, err := Pool.Query(ctx, query)

	var (
		empty = utils.EmptySlice[models.Stock]()
		res   = make([]models.Stock, 0, constants.NumOfStocks)
	)

	if err != nil {
		return empty, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		stock := models.Stock{Segment: "CASH", Exchange: "NSE"}

		err := rows.Scan(&stock.Symbol)

		if err != nil {
			return empty, fmt.Errorf("scanning failed: %w", err)
		}
		stock.Name = stock.Symbol

		res = append(res, stock)
	}

	return res, nil
}

// DeleteDelistedStocks deletes the stock data which is no longer listed on NSE
// Only to be executed on successful completion of the analysis
func DeleteDelistedStocks() {
	log.Println("finding delisted stocks for deletion")
	ctx := context.Background()

	rows, err := Pool.Query(ctx, `
		WITH
			ts_info AS (
				SELECT MAX(timestamp) AS ts
				FROM stock_prices
			),
			delisted AS (
				SELECT symbol
				FROM stock_prices
				GROUP BY symbol
				HAVING MAX(timestamp) <> (
					SELECT ts
					FROM ts_info
				)
			)
		DELETE FROM stock_prices
		WHERE symbol IN (
			SELECT symbol
			FROM delisted
		)
	`)
	if err != nil {
		log.Printf("deletion of delisted stocks failed: %v\n", err)
		return
	}
	defer rows.Close()

	log.Printf("deletion of delisted stocks done")
}
