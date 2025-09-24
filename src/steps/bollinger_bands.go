package steps

import (
	"eeye/src/models"
	"log"
	"math"
)

const (
	MinPoints     = 22
	Period        = 20
	K             = 2
	FlatThreshold = 0.0001
)

func LowerBollingerBandFlatOrVShape(
	strategy string,
	stock *models.Stock,
) func() bool {
	return func() bool {
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
				std_dev := math.Sqrt(variance / float64(Period))
				lbb = append(lbb, avg-K*std_dev)
				sum -= candles[i+1-Period].Close
			}
		}

		lastCandle := candles[length-1]
		isBullishCandle := lastCandle.Open >= lastCandle.Close
		hasCandleFullyCrossedSMA := (isBullishCandle && lastCandle.Open > sma) ||
			(lastCandle.Close > sma)

		if !hasCandleFullyCrossedSMA {
			lbb_len := len(lbb)
			var (
				x1 = 0.0
				y1 = lbb[lbb_len-3]
				x2 = 1.0
				y2 = lbb[lbb_len-2]
				x3 = 2.0
				y3 = lbb[lbb_len-1]
			)

			slope1 := (y2 - y1) / (x2 - x1)
			slope2 := (y3 - y2) / (x3 - x2)

			// check if flat, float div never give perfect 0 there will be some error
			isFlat := math.Abs(slope1) < FlatThreshold && math.Abs(slope2) < FlatThreshold
			isVShape := slope1 < 0 && slope2 > 0
			if isFlat || isVShape {
				return true
			}
		}

		return false
	}
}
