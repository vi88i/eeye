package steps

import (
	"eeye/src/models"
	"eeye/src/utils"
	"log"
)

func computeEMA(candles []models.Candle, period int) []float64 {
	var (
		empty  = utils.EmptySlice[float64]()
		length = len(candles)
	)

	if length < period {
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

// EMAFakeBreakdown creates a function that screens for stocks showing a fake breakdown
// pattern relative to their Exponential Moving Average (EMA). This pattern occurs when
// price temporarily breaks below the EMA but quickly recovers, indicating a false bearish signal.
//
// Parameters:
//   - strategyName: identifier for logging purposes
//   - stock: the stock to analyze
//   - period: the EMA period to use (e.g., 20 for 20-day EMA)
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
			values       = computeEMA(candles, period)
			emaLength    = len(values)
			candleLength = len(candles)
		)

		if emaLength < MinEMAPoints {
			log.Printf("insufficient candles for %v screener: %v\n", strategyName, stock.Symbol)
			return false
		}

		return (candles[candleLength-1].Low <= values[emaLength-1]) &&
			(candles[candleLength-1].High > values[emaLength-1])
	}
}
