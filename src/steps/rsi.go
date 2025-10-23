package steps

import (
	"eeye/src/models"
	"eeye/src/store"
	"eeye/src/utils"
	"log"
	"math"
)

// Rsi screens stocks based on Relative Strength Index (RSI).
// RSI measures momentum on a 0-100 scale to identify overbought/oversold conditions:
//   - Above 70: Overbought (potential sell signal)
//   - Below 30: Oversold (potential buy signal)
//   - 40-60: Neutral/swing zone
type Rsi struct {
	models.StepBaseImpl
	// Test receives RSI values to determine if the stock meets criteria.
	// Parameters:
	//   - rsiValues: Calculated RSI values (14-period standard)
	// Returns true if the stock passes the screening test.
	Test func(rsiValues []float64) bool
}

//revive:disable-next-line exported
func (r *Rsi) Name() string {
	return "RSI screener"
}

//revive:disable-next-line exported
func (r *Rsi) Screen(strategy string, stock *models.Stock) bool {
	const (
		Period = 14 // Standard RSI period (Wilder's original specification)
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

// ComputeRsi calculates the Relative Strength Index using Wilder's smoothing method.
// Algorithm:
//  1. Calculate initial average gain and loss over 'period' candles
//  2. Apply smoothed moving average: avgGain = (prevAvg*(period-1) + currentGain) / period
//  3. Calculate RS (Relative Strength) = avgGain / avgLoss
//  4. Calculate RSI = 100 - (100 / (1 + RS))
//
// Parameters:
//   - candles: Historical price data
//   - period: Lookback period (standard is 14)
//
// Returns:
//   - Slice of RSI values (empty if insufficient data)
func ComputeRsi(candles []models.Candle, period int) []float64 {
	var (
		empty  = utils.EmptySlice[float64]()
		length = len(candles)
	)

	if length < period+1 {
		return empty
	}

	// Calculate initial average gain and loss
	var (
		gain = 0.0
		loss = 0.0
	)
	for i := 1; i <= period; i++ {
		diff := candles[i].Close - candles[i-1].Close
		gain += math.Max(diff, 0)  // Sum positive price changes
		loss += math.Max(-diff, 0) // Sum negative price changes (as positive)
	}

	var (
		p       = float64(period)
		avgGain = gain / p
		avgLoss = loss / p
		values  = make([]float64, 0, length-period)
	)

	// Calculate RSI for each subsequent candle using Wilder's smoothing
	for i := period + 1; i < length; i++ {
		diff := candles[i].Close - candles[i-1].Close
		gain = math.Max(diff, 0)
		loss = math.Max(-diff, 0)

		// Apply Wilder's smoothing formula
		avgGain = ((avgGain * (p - 1)) + gain) / p
		avgLoss = ((avgLoss * (p - 1)) + loss) / p

		// Calculate RSI
		var (
			rs  = avgGain / avgLoss
			rsi = 100 - (100 / (1 + rs))
		)
		values = append(values, rsi)
	}

	return values
}
