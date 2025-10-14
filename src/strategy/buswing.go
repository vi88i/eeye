// Package strategy implements high-level trading strategies screener.
package strategy

import (
	"eeye/src/models"
	"eeye/src/steps"
	"eeye/src/utils"
)

// BullishSwing strategy
type BullishSwing struct {
	models.StrategyBaseImpl
}

//revive:disable-next-line exported
func (b *BullishSwing) Name() string {
	return "Bullish Swing"
}

//revive:disable-next-line exported
func (b *BullishSwing) Execute(stock *models.Stock) {
	var (
		strategyName = b.Name()
		sink         = b.GetSink()
	)

	screeners := []models.Step{
		&steps.BullishCandle{},
		&steps.Volume{
			Test: func(currentVolume float64, averageVolume float64) bool {
				return currentVolume >= averageVolume
			},
		},
		&steps.Rsi{
			Test: func(rsi []float64) bool {
				var (
					length = len(rsi)
					v      = rsi[length-1]
				)

				return v >= 40.0 && v <= 60.0
			},
		},
		&steps.BollingerBands{
			Test: func(candles []models.Candle, sma []float64, lbb []float64, _ []float64) bool {
				return utils.LowerBollingerBandFlatOrVShape(candles, sma, lbb)
			},
		},
	}

	if steps.Execute(strategyName, stock, screeners) {
		sink <- stock
	}
}
