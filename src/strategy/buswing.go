// Package strategy implements high-level trading strategies.
// It combines various technical analysis steps to create complete
// trading strategies and manages their execution.
package strategy

import (
	"eeye/src/models"
	"eeye/src/steps"
)

// BullishSwing strategy
type BullishSwing struct {
	models.StrategyBaseImpl
}

//nolint:revive
func (b *BullishSwing) Name() string {
	return "Bullish Swing"
}

//nolint:revive
func (b *BullishSwing) Execute(stock *models.Stock) {
	var (
		strategyName = b.Name()
		sink         = b.GetSink()
	)

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
			func(rsi []float64) bool {
				var (
					length = len(rsi)
					v      = rsi[length-1]
				)

				return v >= 40.0 && v <= 60.0
			},
		),
		steps.LowerBollingerBandFlatOrVShape(
			strategyName,
			stock,
		),
	}

	if steps.Screen(screeners) {
		sink <- stock
	}
}
