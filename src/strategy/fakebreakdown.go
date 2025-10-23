package strategy

import (
	"eeye/src/models"
	"eeye/src/steps"
	"eeye/src/utils"
	"log"
)

// FakeBreakdown strategy identifies stocks that have broken below support/liquidity levels
// but show bullish reversal patterns, suggesting a potential fake breakdown.
// This strategy looks for:
// 1. A bullish candle pattern (indicating potential reversal)
// 2. Price breaking below a liquidity level but closing above it (fake breakdown)
type FakeBreakdown struct {
	models.StrategyBaseImpl
	Window    int     // Lookback window for identifying liquidity levels
	Tolerance float64 // Price tolerance for clustering levels (e.g., 0.02 for 2%)
	Strength  int     // Minimum touches required for a level to be significant
}

//revive:disable-next-line exported
func (f *FakeBreakdown) Name() string {
	return "Fake Breakdown"
}

//revive:disable-next-line exported
func (f *FakeBreakdown) Execute(stock *models.Stock) {
	var (
		strategyName = f.Name()
		sink         = f.GetSink()
	)

	screeners := []models.Step{
		&steps.BullishCandle{},
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
						log.Printf("[%v] %v did fake breakdown at level %v", strategyName, stock.Symbol, level)
						return currentLow < level && currentClose > level
					},
				)

				return len(fakeBreakdownLevels) > 0
			},
		},
		&steps.Volume{
			Test: func(currentVolume float64, averageVolume float64) bool {
				// Higher volume on the reversal candle adds conviction to the fake breakdown
				// Volume should be above average to confirm strong participation
				return currentVolume >= averageVolume
			},
		},
	}

	if steps.Execute(strategyName, stock, screeners) {
		sink <- stock
	}
}
