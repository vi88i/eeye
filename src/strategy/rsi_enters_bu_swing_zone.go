package strategy

import (
	"eeye/src/models"
	"eeye/src/steps"
	"log"
)

func rsiEntersBullishSwingZoneWorker(strategyName string, in, out chan *models.Stock) {
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
			steps.RSIScreener(
				strategyName,
				stock,
				func(rsi []float64) bool {
					length := len(rsi)
					if length < 2 {
						return false
					}

					var (
						cur  = rsi[length-1]
						prev = rsi[length-2]
					)
					return cur >= 40.0 && prev <= 40.0
				},
			),
		}

		if steps.Executor(screeners) {
			out <- stock
		}

		steps.PurgeCache(stock)
	}
}

func rsiEntersBullishSwingZone(stocks []models.Stock) string {
	return steps.Worker(
		"Bullish Swing Zone RSI",
		stocks,
		rsiEntersBullishSwingZoneWorker,
	)
}
