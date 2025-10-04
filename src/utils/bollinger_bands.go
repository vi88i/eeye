package utils

import (
	"eeye/src/models"
	"math"
)

// LowerBollingerBandFlatOrVShape creates a function that screens for stocks showing
// a flattening or V-shaped pattern in their lower Bollinger Band. This pattern
// often indicates potential trend reversal points.
func LowerBollingerBandFlatOrVShape(candles []models.Candle, sma []float64, lbb []float64) bool {
	const (
		FlatThreshold = 0.0001
	)
	last := len(sma) - 1
	lastCandle := candles[len(candles)-1]
	isBullishCandle := lastCandle.Open >= lastCandle.Close
	hasCandleFullyCrossedSMA := (isBullishCandle && lastCandle.Open > sma[last]) ||
		(lastCandle.Close > sma[last])

	if !hasCandleFullyCrossedSMA {
		lbbLen := len(lbb)
		var (
			x1 = 0.0
			y1 = lbb[lbbLen-3]
			x2 = 1.0
			y2 = lbb[lbbLen-2]
			x3 = 2.0
			y3 = lbb[lbbLen-1]
		)

		slope1 := (y2 - y1) / (x2 - x1)
		slope2 := (y3 - y2) / (x3 - x2)

		isFlat := math.Abs(slope1) < FlatThreshold && math.Abs(slope2) < FlatThreshold
		isVShape := slope1 < 0 && slope2 > 0
		if isFlat || isVShape {
			return true
		}
	}

	return false
}
