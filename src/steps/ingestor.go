package steps

import (
	"eeye/src/api"
	"eeye/src/config"
	"eeye/src/constants"
	"eeye/src/db"
	"eeye/src/models"
	"eeye/src/utils"
	"fmt"
	"log"
	"sync"
	"time"
)

// backFillCandles fetches and stores new candle data for a stock starting from the day
// after the latest candle present in the database up to the current day. It ensures
// that only new data is fetched to avoid duplicates and minimize API calls.
func backFillCandles(stock *models.Stock) error {
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

// ingestionWorker processes stocks from the input channel and backfills their candle data.
func ingestionWorker(in <-chan *models.Stock) {
	for stock := range in {
		if err := backFillCandles(stock); err != nil {
			log.Printf("ingestion failed for %v: %v\n", stock.Symbol, err)
		}
	}
}

// Ingestor updates the historical price data for a stock by fetching new candles
// from the API and storing them in the database. It only fetches data newer than
// the most recent candle in the database to avoid duplicates and minimize API calls.
func Ingestor(stocks []models.Stock) {
	var (
		in = make(chan *models.Stock, constants.IngestionBufferSize)
		wg = sync.WaitGroup{}
	)

	for range constants.NumOfIngestionWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ingestionWorker(in)
		}()
	}

	start := time.Now()
	for i := range stocks {
		in <- &stocks[i]
		if (i+1)%config.TradingAPIConfig.RequestPerSecond == 0 {
			elapsed := time.Since(start)

			if rem := time.Second - elapsed; rem > 0 {
				time.Sleep(rem)
			}

			start = time.Now()
		}
	}
	close(in)

	wg.Wait()
}
