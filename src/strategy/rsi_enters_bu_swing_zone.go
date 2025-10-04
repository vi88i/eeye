package strategy

import (
	"eeye/src/models"
	"eeye/src/steps"
	"log"
)

// RSIEntersBullishSwingZone strategy which detects stocks whose price moved,
// from below baseLine to above baseLine but bounded by upperBound
type RSIEntersBullishSwingZone struct {
	models.StrategyBaseImpl

	baseLine   float64
	upperBound float64
}

//revive:disable-next-line exported
func (r *RSIEntersBullishSwingZone) Name() string {
	return "RSI Enters Bullish Swing Zone"
}

//revive:disable-next-line exported
func (r *RSIEntersBullishSwingZone) Execute(stock *models.Stock) {
	var (
		strategyName = r.Name()
		sink         = r.GetSink()
	)

	if r.baseLine == 0 {
		log.Printf("[%v] baseLine cannot be zero\n", strategyName)
		return
	}

	if r.upperBound == 0 {
		log.Printf("[%v] upperBound cannot be zero\n", strategyName)
		return
	}

	if r.baseLine > r.upperBound {
		log.Printf("[%v] baseLine > upperBound\n", strategyName)
		return
	}

	screeners := []models.Step{
		&steps.BullishCandle{},
		&steps.RSI{
			Test: func(rsi []float64) bool {
				length := len(rsi)
				if length < 2 {
					return false
				}

				var (
					cur  = rsi[length-1]
					prev = rsi[length-2]
				)
				return cur >= r.baseLine && prev <= r.baseLine && cur <= r.upperBound && cur > prev
			},
		},
	}

	if steps.Execute(strategyName, stock, screeners) {
		sink <- stock
	}
}
