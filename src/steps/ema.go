package steps

import (
	"eeye/src/models"
	"eeye/src/store"
	"eeye/src/utils"
	"fmt"
	"log"
)

// EMA creates a function that screens for stocks based on their Exponential Moving Average (EMA).
type EMA struct {
	Period int
	Test   func(candles []models.Candle, emas []float64) bool
}

//revive:disable-next-line exported
func (e *EMA) Name() string {
	return fmt.Sprintf("EMA %v screener", e.Period)
}

//revive:disable-next-line exported
func (e *EMA) Screen(strategy string, stock *models.Stock) bool {
	const (
		MinEMAPoints = 1
	)

	step := e.Name()

	candles, err := store.Get(stock)
	if err != nil {
		return false
	}

	var (
		values    = ComputeEMA(candles, e.Period)
		emaLength = len(values)
	)

	if emaLength < MinEMAPoints {
		log.Printf("[%v - %v] insufficient candles: %v\n", strategy, step, stock.Symbol)
		return false
	}

	test := e.Test(candles, values)
	if !test {
		log.Printf("[%v - %v] test failed: %v\n", strategy, step, stock.Symbol)
	}
	return test
}

// ComputeEMA is helper method to compute the EMA values
func ComputeEMA(candles []models.Candle, period int) []float64 {
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
