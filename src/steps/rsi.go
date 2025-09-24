package steps

import (
	"eeye/src/models"
	"log"
	"math"
)

const (
	n = 14
)

func RSIScreener(
	strategy string,
	stock *models.Stock,
	screen func(currentRSI float64) bool,
) func() bool {
	return func() bool {
		candles, err := getCachedCandles(stock)
		if err != nil {
			return false
		}

		length := len(candles)
		/*
			We need 15 candles minimum to get started with RSI
		*/
		if length < n+1 {
			log.Printf("insufficient candles for RSI screener: %v\n", stock.Symbol)
			return false
		}

		var (
			gain     = 0.0
			loss     = 0.0
			avg_gain = 0.0
			avg_loss = 0.0
		)
		for i := 1; i <= n; i++ {
			diff := candles[i].Close - candles[i-1].Close
			gain += math.Max(diff, 0)
			loss += math.Max(-diff, 0)
		}

		avg_gain = gain / n
		avg_loss = loss / n

		for i := n + 1; i < length; i++ {
			diff := candles[i].Close - candles[i-1].Close
			gain = math.Max(diff, 0)
			loss = math.Max(-diff, 0)

			avg_gain = ((avg_gain * (n - 1)) + gain) / n
			avg_loss = ((avg_loss * (n - 1)) + loss) / n
		}

		rs := avg_gain / avg_loss
		rsi := 100 - (100 / (1 + rs))
		test := screen(rsi)
		if !test {
			log.Printf("[%v] %v failed RSI test\n", strategy, stock.Symbol)
		}
		return test
	}
}
