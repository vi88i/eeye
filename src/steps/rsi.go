package steps

import (
	"eeye/src/models"
	"eeye/src/utils"
	"log"
	"math"
)

func computeRSI(candles []models.Candle, period int) []float64 {
	var (
		empty  = utils.EmptySlice[float64]()
		length = len(candles)
	)

	if length < period+1 {
		return empty
	}

	var (
		gain = 0.0
		loss = 0.0
	)
	for i := 1; i <= period; i++ {
		diff := candles[i].Close - candles[i-1].Close
		gain += math.Max(diff, 0)
		loss += math.Max(-diff, 0)
	}

	var (
		p       = float64(period)
		avgGain = gain / p
		avgLoss = loss / p
		values  = make([]float64, 0, length-period)
	)

	for i := period + 1; i < length; i++ {
		diff := candles[i].Close - candles[i-1].Close
		gain = math.Max(diff, 0)
		loss = math.Max(-diff, 0)

		avgGain = ((avgGain * (p - 1)) + gain) / p
		avgLoss = ((avgLoss * (p - 1)) + loss) / p

		var (
			rs  = avgGain / avgLoss
			rsi = 100 - (100 / (1 + rs))
		)
		values = append(values, rsi)
	}

	return values
}

// RSIScreener creates a function that screens stocks based on their Relative Strength
// Index (RSI) value. It calculates the 14-period RSI and applies a custom screening
// function to determine if the stock meets the criteria.
//
// Parameters:
//   - strategy: identifier for logging purposes
//   - stock: the stock to analyze
//   - screen: a function that takes the current RSI value and returns true if the
//     stock passes the screening criteria
func RSIScreener(
	strategy string,
	stock *models.Stock,
	screen func(rsiValues []float64) bool,
) func() bool {
	return func() bool {
		const (
			Period = 14
		)

		candles, err := getCachedCandles(stock)
		if err != nil {
			return false
		}

		var (
			rsi       = computeRSI(candles, Period)
			rsiLength = len(rsi)
		)

		if rsiLength == 0 {
			log.Printf("insufficient candles for RSI: %v\n", stock.Symbol)
			return false
		}

		test := screen(rsi)
		if !test {
			log.Printf("%v failed %v RSI test\n", stock.Symbol, strategy)
		}
		return test
	}
}
