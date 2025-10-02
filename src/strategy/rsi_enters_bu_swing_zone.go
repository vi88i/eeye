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

//nolint:revive
func (r *RSIEntersBullishSwingZone) Name() string {
	return "RSI Enters Bullish Swing Zone"
}

//nolint:revive
func (r *RSIEntersBullishSwingZone) Execute(stock *models.Stock) {
	var (
		strategyName = r.Name()
		sink         = r.GetSink()
	)

	if r.baseLine == 0 {
		log.Println("baseLine cannot be zero")
		return
	}

	if r.upperBound == 0 {
		log.Println("upperBound cannot be zero")
		return
	}

	if r.baseLine > r.upperBound {
		log.Println("baseLine > upperBound")
		return
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
				return cur >= r.baseLine && prev <= r.baseLine && cur <= r.upperBound && cur > prev
			},
		),
	}

	if steps.Screen(screeners) {
		sink <- stock
	}
}
