package steps

import (
	"eeye/src/models"
	"log"
	"math"
)

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
	screen func(currentRSI float64) bool,
) func() bool {
	return func() bool {
		const (
			Period = 14
		)

		candles, err := getCachedCandles(stock)
		if err != nil {
			return false
		}

		length := len(candles)
		if length < Period+1 {
			log.Printf("insufficient candles for RSI screener: %v\n", stock.Symbol)
			return false
		}

		var (
			gain = 0.0
			loss = 0.0
		)
		for i := 1; i <= Period; i++ {
			diff := candles[i].Close - candles[i-1].Close
			gain += math.Max(diff, 0)
			loss += math.Max(-diff, 0)
		}

		// Initialize average gain and loss
		avgGain := gain / Period
		avgLoss := loss / Period

		for i := Period + 1; i < length; i++ {
			diff := candles[i].Close - candles[i-1].Close
			gain = math.Max(diff, 0)
			loss = math.Max(-diff, 0)

			avgGain = ((avgGain * (Period - 1)) + gain) / Period
			avgLoss = ((avgLoss * (Period - 1)) + loss) / Period
		}

		rs := avgGain / avgLoss
		rsi := 100 - (100 / (1 + rs))
		test := screen(rsi)
		if !test {
			log.Printf("[%v] %v failed RSI test\n", strategy, stock.Symbol)
		}
		return test
	}
}
