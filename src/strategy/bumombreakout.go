package strategy

import (
	"eeye/src/models"
	"eeye/src/steps"
	"eeye/src/utils"
	"math"
)

// BullishMomentumBreakout helps to identify stocks which could potentially breakout,
// and go in momentum
type BullishMomentumBreakout struct {
	models.StrategyBaseImpl
}

//revive:disable-next-line exported
func (b *BullishMomentumBreakout) Name() string {
	return "Bullish momentum"
}

//revive:disable-next-line exported
func (b *BullishMomentumBreakout) Execute(stock *models.Stock) {
	var (
		strategyName = b.Name()
		sink         = b.GetSink()
	)

	screeners := []models.Step{
		&steps.BullishCandle{},
		&steps.Rsi{
			Test: func(rsi []float64) bool {
				length := len(rsi)

				if length < 2 {
					return false
				}

				return rsi[length-2] <= 60 &&
					rsi[length-1] > 60
			},
		},
		&steps.Ema{
			Period: 50,
			Test: func(candles []models.Candle, emas []float64) bool {
				var (
					emaLength    = len(emas)
					candleLength = len(candles)
				)

				return candles[candleLength-1].Close > emas[emaLength-1]
			},
		},
		&steps.BollingerBands{
			Test: func(candles []models.Candle, _, _, ubb []float64) bool {
				var (
					ubbLength    = len(ubb)
					candleLength = len(candles)
				)

				return candles[candleLength-1].High > ubb[ubbLength-1]
			},
		},
		&steps.EmaCrossover{
			Periods: []int{5, 13, 26, 50, 200},
			Test: func(emas [][]float64) bool {
				prev := math.Inf(1)
				for i := range emas {
					next := utils.Last(emas[i], math.Inf(1))
					if prev < next {
						return false
					}
					prev = next
				}

				return true
			},
		},
	}

	if steps.Execute(strategyName, stock, screeners) {
		sink <- stock
	}
}
