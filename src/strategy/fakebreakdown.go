package strategy

import (
	"eeye/src/models"
	"eeye/src/steps"
	"eeye/src/utils"
	"log"
)

// FakeBreakdown identifies stocks that have broken below support/liquidity levels
// but show bullish reversal patterns, suggesting a potential fake breakdown (also known as a bull trap fake-out).
// This strategy looks for:
//   - A bullish candle pattern (indicating potential reversal)
//   - Price breaking below a support/liquidity level but closing above it (the fake breakdown)
//   - Above-average volume confirming the reversal
//
// Trading Logic:
//   - Support levels are identified using a clustering algorithm on historical lows
//   - A fake breakdown occurs when bears push price below support (low < level)
//     but bulls regain control and close above it (close > level)
//   - This "bear trap" often leads to a strong move higher as bears cover shorts
//   - Volume confirmation ensures the reversal has conviction
//
// Example Scenario:
//   - Stock has support at $100 (tested multiple times)
//   - Current candle: low = $99.50, close = $100.50
//   - This is a fake breakdown - bears failed to sustain the break
//
// Ideal For: Swing traders looking to capitalize on failed breakdowns and reversals
// Timeframe: Works on all timeframes, most effective on daily charts
// Risk Profile: Medium - requires both technical pattern and volume confirmation
type FakeBreakdown struct {
	models.StrategyBaseImpl
	Window    int     // Lookback window for identifying liquidity levels (e.g., 5 for recent levels)
	Tolerance float64 // Price tolerance for clustering levels (e.g., 0.01 for 1%, 0.02 for 2%)
	Strength  int     // Minimum touches required for a level to be significant (e.g., 3)
}

// Name returns the strategy identifier.
//
// Returns:
//   - The name of this strategy ("Fake Breakdown")
//
//revive:disable-next-line exported
func (f *FakeBreakdown) Name() string {
	return "Fake Breakdown"
}

// Execute runs the FakeBreakdown strategy on the given stock.
// It applies three screening steps:
//  1. BullishCandle: Confirms bullish reversal pattern
//  2. LiquidityLevels: Identifies support levels and checks for fake breakdown pattern
//  3. Volume: Ensures above-average volume for conviction
//
// If all conditions are met, the stock is sent to the strategy's output sink.
//
// Parameters:
//   - stock: The stock to analyze for fake breakdown pattern
//
//revive:disable-next-line exported
func (f *FakeBreakdown) Execute(stock *models.Stock) {
	var (
		strategyName = f.Name()
		sink         = f.GetSink()
	)

	screeners := []models.Step{
		// Step 1: Confirm bullish candlestick pattern showing reversal
		&steps.BullishCandle{},

		// Step 2: Check for fake breakdown at support levels
		&steps.LiquidityLevels{
			Window:    f.Window,
			Tolerance: f.Tolerance,
			Strength:  f.Strength,
			Test: func(candles []models.Candle, supports []float64, _ []float64) bool {
				// A fake breakdown occurs when:
				// - Price has broken below a support level (low is below the level)
				// - But closes above it (close is above the level)
				// This suggests bears tried to push price down but bulls regained control

				if len(supports) == 0 || len(candles) == 0 {
					return false
				}

				// Get the current (latest) candle using utils.Last helper
				currentCandle := utils.Last(candles, models.Candle{})
				currentLow := currentCandle.Low
				currentClose := currentCandle.Close

				// Check if price created a fake breakdown against any liquidity level
				// Using functional approach to check if any level satisfies the condition
				fakeBreakdownLevels := utils.Filter(
					supports,
					func(level float64, _ int) bool {
						// Fake breakdown condition:
						// Low went below the level BUT close is above the level
						return currentLow < level && currentClose > level
					},
				)

				log.Printf("[%v] %v did fake breakdown at levels %v", strategyName, stock.Symbol, fakeBreakdownLevels)
				return len(fakeBreakdownLevels) > 0
			},
		},

		// Step 3: Verify volume is above average
		&steps.Volume{
			Test: func(currentVolume float64, averageVolume float64) bool {
				// Higher volume on the reversal candle adds conviction to the fake breakdown
				// Volume should be at or above average to confirm strong participation
				// This prevents false signals from low-volume, low-conviction reversals
				return currentVolume >= averageVolume
			},
		},
	}

	// Execute all screening steps; if all three pass, send stock to output sink
	if steps.Execute(strategyName, stock, screeners) {
		sink <- stock
	}
}
