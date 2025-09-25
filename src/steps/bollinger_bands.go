// Package steps implements various technical analysis indicators and trading strategies.
// It provides functions for analyzing price patterns, calculating technical indicators,
// and screening stocks based on specific criteria.
package steps

import (
	"eeye/src/models"
	"log"
	"math"
)

// LowerBollingerBandFlatOrVShape creates a function that screens for stocks showing
// a flattening or V-shaped pattern in their lower Bollinger Band. This pattern
// often indicates potential trend reversal points.
//
// The function uses a 20-period Bollinger Band with 2 standard deviations and
// requires at least 22 data points for calculation.
func LowerBollingerBandFlatOrVShape(
	strategy string,
	stock *models.Stock,
) func() bool {
	return func() bool {
		const (
			MinPoints     = 22
			Period        = 20
			K             = 2
			FlatThreshold = 0.0001
		)

		candles, err := getCachedCandles(stock)
		if err != nil {
			return false
		}

		length := len(candles)
		if length < MinPoints {
			log.Printf(
				"insufficient candles for %v lower bollinger band signal: %v\n",
				strategy,
				stock.Symbol,
			)
			return false
		}

		sum := 0.0
		sma := 0.0
		lbb := make([]float64, 0, len(candles)-Period+1)
		for i := range candles {
			sum += candles[i].Close
			if i+1 >= Period {
				avg := sum / float64(Period)
				sma = avg
				variance := 0.0
				for j := i + 1 - Period; j <= i; j++ {
					diff := candles[j].Close - avg
					variance = variance + diff*diff
				}
				stdDev := math.Sqrt(variance / float64(Period))
				lbb = append(lbb, avg-K*stdDev)
				sum -= candles[i+1-Period].Close
			}
		}

		lastCandle := candles[length-1]
		isBullishCandle := lastCandle.Open >= lastCandle.Close
		hasCandleFullyCrossedSMA := (isBullishCandle && lastCandle.Open > sma) ||
			(lastCandle.Close > sma)

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
}
