package strategy

import (
	"eeye/src/models"
	"eeye/src/steps"
	"eeye/src/utils"
)

// LowerBollingerBandBullish strategy which checks if lower bollinger band,
// does a 'V' shape or have a positive slope
type LowerBollingerBandBullish struct {
	models.StrategyBaseImpl
}

//revive:disable-next-line exported
func (l *LowerBollingerBandBullish) Name() string {
	return "Lower Bollinger Band Bullish"
}

//revive:disable-next-line exported
func (l *LowerBollingerBandBullish) Execute(stock *models.Stock) {
	var (
		strategyName = l.Name()
		sink         = l.GetSink()
	)

	screeners := []models.Step{
		&steps.BullishCandle{},
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
