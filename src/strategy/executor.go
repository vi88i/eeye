package strategy

import (
	"eeye/src/models"
	"log"
)

// Executor runs all trading strategies on the given list of stocks.
// It executes multiple strategies in parallel and logs their results.
// Current strategies include:
//   - Lower Bollinger Band bullish pattern
//   - Bullish swing pattern
//   - EMA fake breakdown pattern
func Executor(stocks []models.Stock) {
	results := []string{
		lowerBollingerBandBullish(stocks),
		bullishSwing(stocks),
		emaFakeBreakdown(stocks, 50),
	}

	log.Println("================= Strategy Results =================")
	for _, result := range results {
		log.Println(result)
	}
}
