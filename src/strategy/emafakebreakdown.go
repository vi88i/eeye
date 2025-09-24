package strategy

import (
	"eeye/src/models"
	"eeye/src/steps"
	"fmt"
	"log"
)

func emaFakeBreakdownWorker(period int) func(strategyName string, in, out chan *models.Stock) {
	return func(strategyName string, in, out chan *models.Stock) {
		for stock := range in {
			if err := steps.Ingestor(stock); err != nil {
				log.Printf("ingestion failed for %v: %v\n", stock.Symbol, err)
				continue
			}

			if err := steps.Extractor(stock); err != nil {
				log.Printf("historical data extraction failed for %v: %v", stock.Symbol, err)
				continue
			}

			screeners := []func() bool{
				steps.BullishCandleScreener(
					strategyName,
					stock,
				),
				steps.EMAFakeBreakdown(
					strategyName,
					stock,
					50,
				),
			}

			if steps.Executor(screeners) {
				out <- stock
			}

			steps.PurgeCache(stock)
		}
	}
}

func emaFakeBreakdown(stocks []models.Stock, period int) string {
	return steps.Worker(
		fmt.Sprintf("EMA %v fake breakdown", period),
		stocks,
		emaFakeBreakdownWorker(period),
	)
}
