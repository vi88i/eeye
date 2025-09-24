package strategy

import (
	"eeye/src/models"
	"eeye/src/steps"
	"log"
)

func bullishSwingWorker(strategyName string, in, out chan *models.Stock) {
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
			steps.VolumeScreener(
				strategyName,
				stock,
				func(currentVolume float64, averageVolume float64) bool {
					return currentVolume >= averageVolume
				},
			),
			steps.RSIScreener(
				strategyName,
				stock,
				func(currentRSI float64) bool {
					return currentRSI >= 40.0 && currentRSI <= 60.0
				},
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

func bullishSwing(stocks []models.Stock) string {
	return steps.Worker(
		"Bullish Swing",
		stocks,
		bullishSwingWorker,
	)
}
