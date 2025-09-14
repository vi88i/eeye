package steps

import (
	"eeye/src/api"
	"eeye/src/db"
	"eeye/src/models"
	"eeye/src/utils"
	"fmt"
	"time"
)

// Ingestor fetches latest candles for all stocks and backfills the database.
func Ingestor(stock *models.Stock) error {
	latestCandle, err := db.GetLastCandle(stock.Symbol)
	if err != nil {
		return fmt.Errorf("failed to fetch latest candle for %v: %w", stock.Symbol, err)
	}

	start := utils.GetFormattedTimestamp(latestCandle.Timestamp.AddDate(0, 0, 1))
	end := utils.GetFormattedTimestamp(time.Now().UTC().Truncate(24 * time.Hour).AddDate(0, 0, 1))
	newCandles, err := api.GetCandles(stock, start, end)
	if err != nil {
		return fmt.Errorf("failed to fetch latest candles for %v: %w", stock.Symbol, err)
	}

	if err = db.BackfillCandles(stock, newCandles); err != nil {
		return fmt.Errorf("failed to ingest data for %v: %w", stock.Symbol, err)
	}

	return nil
}
