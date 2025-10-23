package steps

import (
	"eeye/src/models"
	"eeye/src/store"
	"eeye/src/utils"
	"fmt"
	"log"
)

// Ema screens stocks based on Exponential Moving Average (EMA) analysis.
// EMA gives more weight to recent prices, making it more responsive to price changes than SMA.
type Ema struct {
	models.StepBaseImpl
	// Period is the number of candles to use for EMA calculation
	Period int
	// Test is a custom function that receives candles and EMA values to determine screening criteria.
	// Parameters:
	//   - candles: Historical price data
	//   - emas: Calculated EMA values for the given period
	// Returns true if the stock passes the screening test.
	Test func(candles []models.Candle, emas []float64) bool
}

//revive:disable-next-line exported
func (e *Ema) Name() string {
	return fmt.Sprintf("EMA %v screener", e.Period)
}

//revive:disable-next-line exported
func (e *Ema) Screen(strategy string, stock *models.Stock) bool {
	const (
		MinEMAPoints = 1
	)

	step := e.Name()

	candles, err := store.Get(stock)
	if err != nil {
		return false
	}

	var (
		values    = ComputeEma(candles, e.Period)
		emaLength = len(values)
	)

	if emaLength < MinEMAPoints {
		log.Printf("[%v - %v] insufficient candles: %v\n", strategy, step, stock.Symbol)
		return false
	}

	return e.TruthyCheck(
		strategy,
		step,
		stock,
		func() bool {
			return e.Test(candles, values)
		},
	)
}

// ComputeEma calculates the Exponential Moving Average for the given period.
// Algorithm:
//  1. Calculate SMA for first 'period' candles as initial EMA
//  2. Apply EMA formula: EMA = alpha * currentPrice + (1-alpha) * previousEMA
//  3. alpha = 2 / (period + 1) - smoothing factor
//
// Parameters:
//   - candles: Historical price data
//   - period: Number of periods for EMA calculation
//
// Returns:
//   - Slice of EMA values (empty if insufficient data)
func ComputeEma(candles []models.Candle, period int) []float64 {
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
		alpha  = 2.0 / (float64(period) + 1) // Smoothing factor
	)

	// Calculate initial SMA as the first EMA value
	for i := range period {
		sum += candles[i].Close
	}
	values = append(values, sum/float64(period))

	// Apply EMA formula for subsequent values
	for i, j := period, 0; i < length; i, j = i+1, j+1 {
		values = append(values, alpha*candles[i].Close+(1.0-alpha)*values[j])
	}

	return values
}
