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
	`, symbol, config.DBConfig.Tz)

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

func FetchAllCandles(stock *models.Stock) ([]models.Candle, error) {
	log.Printf("fetching all candles: %v\n", stock.Symbol)
	ctx := context.Background()

	rows, err := Pool.Query(ctx, `
		SELECT symbol, open, close, high, low, (timestamp AT TIME ZONE $2) as timestamp, volume
		FROM stock_prices
		WHERE symbol = $1
		ORDER BY timestamp ASC
	`, stock.Symbol, config.DBConfig.Tz)

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
