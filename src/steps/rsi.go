package steps

import (
	"eeye/src/models"
	"eeye/src/store"
	"eeye/src/utils"
	"log"
	"math"
)

// Rsi creates a function that screens stocks based on their Relative Strength
// Index (RSI) value. It calculates the 14-period RSI and applies a custom screening
// function to determine if the stock meets the criteria.
type Rsi struct {
	models.StepBaseImpl
	Test func(rsiValues []float64) bool
}

//revive:disable-next-line exported
func (r *Rsi) Name() string {
	return "RSI screener"
}

//revive:disable-next-line exported
func (r *Rsi) Screen(strategy string, stock *models.Stock) bool {
	const (
		Period = 14
	)

	step := r.Name()

	candles, err := store.Get(stock)
	if err != nil {
		return false
	}

	var (
		rsi       = ComputeRsi(candles, Period)
		rsiLength = len(rsi)
	)

	if rsiLength == 0 {
		log.Printf("[%v - %v] insufficient candles: %v\n", strategy, step, stock.Symbol)
		return false
	}

	return r.TruthyCheck(
		strategy,
		step,
		stock,
		func() bool {
			return r.Test(rsi)
		},
	)
}

// ComputeRsi is helper method to compute the RSI values
func ComputeRsi(candles []models.Candle, period int) []float64 {
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
