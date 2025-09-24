package steps

import (
	"eeye/src/models"
	"eeye/src/utils"
	"log"
)

func ema(candles []models.Candle, stock *models.Stock, period int) []float64 {
	var (
		empty  = utils.EmptySlice[float64]()
		length = len(candles)
	)

	if length < period {
		log.Printf("insufficient candles for EMA %v: %v", period, stock.Symbol)
		return empty
	}

	var (
		values = make([]float64, 0, length-period+1)
		sum    = 0.0
		alpha  = 2.0 / (float64(period) + 1)
	)

	for i := range period {
		sum += candles[i].Close
	}
	values = append(values, sum/float64(period))

	for i, j := period, 0; i < length; i, j = i+1, j+1 {
		values = append(values, alpha*candles[i].Close+(1.0-alpha)*values[j])
	}

	return values
}

func EMAFakeBreakdown(
	strategyName string,
	stock *models.Stock,
	period int,
) func() bool {
	return func() bool {
		const (
			MinEMAPoints = 1
		)

		candles, err := getCachedCandles(stock)
		if err != nil {
			return false
		}

		var (
			values       = ema(candles, stock, period)
			emaLength    = len(values)
			candleLength = len(candles)
		)

		if emaLength < MinEMAPoints {
			log.Printf("insufficient candles for %v screener: %v", strategyName, stock.Symbol)
			return false
		}

		return (candles[candleLength-1].Low <= values[emaLength-1]) &&
			(candles[candleLength-1].High > values[emaLength-1])
	}
}
