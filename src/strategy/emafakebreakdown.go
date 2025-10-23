package strategy

import (
	"eeye/src/models"
	"eeye/src/steps"
	"fmt"
)

// EmaFakeBreakdown identifies stocks that have broken below an Exponential Moving Average (EMA)
// but show immediate bullish reversal, suggesting a fake breakdown or bear trap.
// This strategy looks for:
//   - A bullish candle pattern (indicating reversal)
//   - Price action where low goes below EMA but high is above EMA (wick below, close above)
//
// The EMA acts as dynamic support, and this pattern occurs when:
//  1. Bears push price below the EMA (low <= EMA)
//  2. Bulls regain control and push price back above (high > EMA)
//  3. The candle closes bullish, showing rejection of lower prices
//
// This is a common pattern at key EMA levels (e.g., 50-day, 200-day) and can signal
// the end of a pullback and resumption of the uptrend.
//
// Ideal For: Trend followers looking for pullback entries in uptrends
// Timeframe: Works on all timeframes; commonly used with 50-day or 200-day EMA on daily charts
// Risk Profile: Medium - requires confirmation from both price action and EMA interaction
type EmaFakeBreakdown struct {
	models.StrategyBaseImpl
	period int // The EMA period to use for support detection (e.g., 50, 200)
}

// Name returns the strategy identifier with the EMA period.
//
// Returns:
//   - The name of this strategy (e.g., "EMA 50 fake breakdown")
//
//revive:disable-next-line exported
func (e *EmaFakeBreakdown) Name() string {
	return fmt.Sprintf("EMA %v fake breakdown", e.period)
}

// Execute runs the EmaFakeBreakdown strategy on the given stock.
// It applies two screening steps:
//  1. BullishCandle: Confirms a bullish reversal pattern
//  2. Ema: Checks if price created a fake breakdown below the EMA
//     (low went below but high stayed above, indicating rejection)
//
// If both conditions are met, the stock is sent to the strategy's output sink.
//
// Parameters:
//   - stock: The stock to analyze for EMA fake breakdown pattern
//
//revive:disable-next-line exported
func (e *EmaFakeBreakdown) Execute(stock *models.Stock) {
	var (
		strategyName = e.Name()
		sink         = e.GetSink()
	)

	screeners := []models.Step{
		// Step 1: Confirm bullish candlestick pattern showing reversal
		&steps.BullishCandle{},

		// Step 2: Check for fake breakdown below EMA
		// Condition: candle low breached the EMA but candle high stayed above it
		// This creates a wick below the EMA with body above, showing rejection
		&steps.Ema{
			Period: e.period,
			Test: func(candles []models.Candle, ema []float64) bool {
				candleLength := len(candles)
				emaLength := len(ema)

				// Fake breakdown occurs when:
				// - Low is at or below EMA (bears tested support)
				// - High is above EMA (bulls defended and pushed back)
				return (candles[candleLength-1].Low <= ema[emaLength-1]) &&
					(candles[candleLength-1].High > ema[emaLength-1])
			},
		},
	}

	// Execute all screening steps; if both pass, send stock to output sink
	if steps.Execute(strategyName, stock, screeners) {
		sink <- stock
	}
}
