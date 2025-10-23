package strategy

import (
	"eeye/src/models"
	"eeye/src/steps"
	"log"
)

// RsiEntersBullishSwingZone identifies stocks whose RSI has just entered the bullish swing zone,
// indicating emerging bullish momentum. This strategy looks for:
//   - RSI crossing above a baseline level (e.g., 40) from below
//   - RSI staying within an upper bound (e.g., 60) to avoid overbought conditions
//   - Upward RSI momentum (current RSI > previous RSI)
//   - Bullish candlestick pattern confirming the momentum shift
//
// Trading Logic:
//   - The zone between baseLine and upperBound (typically 40-60) is considered the "sweet spot"
//   - RSI crossing into this zone from below indicates recovering strength without being overbought
//   - This often occurs at the beginning of a new bullish trend or after a healthy pullback
//
// Example Configuration:
//   - baseLine: 40 (entry into strength zone)
//   - upperBound: 60 (not yet overbought)
//
// Ideal For: Swing traders looking to catch momentum shifts early
// Timeframe: Best on daily or 4-hour charts
// Risk Profile: Medium - provides early entry with confirmation from multiple factors
type RsiEntersBullishSwingZone struct {
	models.StrategyBaseImpl

	baseLine   float64 // RSI level that must be crossed from below (e.g., 40)
	upperBound float64 // Maximum RSI level to avoid overbought conditions (e.g., 60)
}

// Name returns the strategy identifier.
//
// Returns:
//   - The name of this strategy ("RSI Enters Bullish Swing Zone")
//
//revive:disable-next-line exported
func (r *RsiEntersBullishSwingZone) Name() string {
	return "RSI Enters Bullish Swing Zone"
}

// Execute runs the RsiEntersBullishSwingZone strategy on the given stock.
// It first validates the configuration parameters, then applies two screening steps:
//  1. BullishCandle: Confirms bullish price action
//  2. Rsi: Checks if RSI has crossed into the swing zone with upward momentum
//
// Validation checks:
//   - baseLine and upperBound must be non-zero
//   - baseLine must be less than upperBound
//
// RSI conditions (all must be true):
//   - Previous RSI was at or below baseLine
//   - Current RSI is at or above baseLine
//   - Current RSI is at or below upperBound
//   - Current RSI is higher than previous RSI (upward momentum)
//
// If all conditions are met, the stock is sent to the strategy's output sink.
//
// Parameters:
//   - stock: The stock to analyze for RSI entry into bullish swing zone
//
//revive:disable-next-line exported
func (r *RsiEntersBullishSwingZone) Execute(stock *models.Stock) {
	var (
		strategyName = r.Name()
		sink         = r.GetSink()
	)

	// Validate configuration parameters
	if r.baseLine == 0 {
		log.Printf("[%v] baseLine cannot be zero\n", strategyName)
		return
	}

	if r.upperBound == 0 {
		log.Printf("[%v] upperBound cannot be zero\n", strategyName)
		return
	}

	if r.baseLine > r.upperBound {
		log.Printf("[%v] baseLine > upperBound\n", strategyName)
		return
	}

	screeners := []models.Step{
		// Step 1: Confirm bullish candlestick pattern
		&steps.BullishCandle{},

		// Step 2: Check if RSI has entered the bullish swing zone
		&steps.Rsi{
			Test: func(rsi []float64) bool {
				length := len(rsi)
				if length < 2 {
					return false
				}

				var (
					cur  = rsi[length-1] // Current RSI value
					prev = rsi[length-2] // Previous RSI value
				)

				// All conditions must be true:
				// 1. Current RSI is at or above baseline (entered the zone)
				// 2. Previous RSI was at or below baseline (crossed from below)
				// 3. Current RSI is within upper bound (not overbought)
				// 4. RSI is rising (current > previous)
				return cur >= r.baseLine && prev <= r.baseLine && cur <= r.upperBound && cur > prev
			},
		},
	}

	// Execute all screening steps; if both pass, send stock to output sink
	if steps.Execute(strategyName, stock, screeners) {
		sink <- stock
	}
}
