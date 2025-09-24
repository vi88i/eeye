package strategy

import (
	"eeye/src/models"
	"eeye/src/steps"
	"log"
)

func llbBullishWorker(strategyName string, in chan *models.Stock, out chan *models.Stock) {
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
			steps.LowerBollingerBandFlatOrVShape(
				strategyName,
				stock,
			),
		}

		if steps.Executor(screeners) {
			out <- stock
		}

		steps.PurgeCache(stock)
	}
}

func lowerBollingerBandBullish(stocks []models.Stock) string {
	return steps.Worker(
		"Lower Bollinger Band Bullish",
		stocks,
		llbBullishWorker,
	)
}
