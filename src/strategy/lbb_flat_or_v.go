package strategy

import (
	"eeye/src/models"
	"eeye/src/steps"
)

// LowerBollingerBandBullish strategy which checks if lower bollinger band,
// does a 'V' shape or have a positive slope
type LowerBollingerBandBullish struct {
	models.StrategyBaseImpl
}

//nolint:revive
func (l *LowerBollingerBandBullish) Name() string {
	return "Lower Bollinger Band Bullish"
}

//nolint:revive
func (l *LowerBollingerBandBullish) Execute(stock *models.Stock) {
	var (
		strategyName = l.Name()
		sink         = l.GetSink()
	)

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

	if steps.Screen(screeners) {
		sink <- stock
	}
}
