package steps

import (
	"eeye/src/api"
	"eeye/src/db"
	"eeye/src/models"
	"eeye/src/utils"
	"fmt"
	"time"
)

func Ingestor(stock *models.Stock) error {
	latestCandle, err := db.GetLastCandle(stock.Symbol)
	if err != nil {
		return fmt.Errorf("failed to fetch latest candle for %v: %w", stock.Symbol, err)
	}

	startOfDayPlusOne := func(t time.Time) time.Time {
		return t.UTC().Truncate(24*time.Hour).AddDate(0, 0, 1)
	}

	var (
		start = utils.GetFormattedTimestamp(startOfDayPlusOne(latestCandle.Timestamp))
		end   = utils.GetFormattedTimestamp(startOfDayPlusOne(time.Now()))
	)

	newCandles, err := api.GetCandles(stock, start, end)
	if err != nil {
		return fmt.Errorf("failed to fetch latest candles for %v: %w", stock.Symbol, err)
	}

	if err = db.BackfillCandles(stock, newCandles); err != nil {
		return fmt.Errorf("failed to ingest data for %v: %w", stock.Symbol, err)
	}

	return nil
}
