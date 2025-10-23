// Package strategy implements high-level trading strategies screener.
package strategy

import (
	"eeye/src/models"
	"eeye/src/steps"
	"eeye/src/utils"
)

// BullishSwing identifies stocks with bullish reversal potential in a swing trading context.
// This strategy combines multiple technical indicators to find stocks that are:
//   - Showing bullish price action (bullish candle patterns)
//   - Trading with healthy volume (at or above average)
//   - In the swing zone according to RSI (40-60 range, neither overbought nor oversold)
//   - Showing support from the lower Bollinger Band (flat or V-shaped pattern)
//
// Ideal For: Swing traders looking for short to medium-term bullish reversals
// Timeframe: Works best on daily charts
// Risk Profile: Medium - requires confirmation from multiple indicators
type BullishSwing struct {
	models.StrategyBaseImpl
}

// Name returns the strategy identifier.
//
// Returns:
//   - The name of this strategy ("Bullish Swing")
//
//revive:disable-next-line exported
func (b *BullishSwing) Name() string {
	return "Bullish Swing"
}

// Execute runs the BullishSwing strategy on the given stock.
// It applies a series of screening steps in sequence:
//  1. BullishCandle: Confirms a bullish candlestick pattern
//  2. Volume: Ensures volume is at or above average (confirmation of interest)
//  3. RSI: Checks if RSI is in the swing zone (40-60) indicating balanced momentum
//  4. BollingerBands: Verifies lower band shows flat or V-shape pattern (support forming)
//
// If all screening steps pass, the stock is sent to the strategy's output sink.
//
// Parameters:
//   - stock: The stock to analyze for bullish swing potential
//
//revive:disable-next-line exported
func (b *BullishSwing) Execute(stock *models.Stock) {
	var (
		strategyName = b.Name()
		sink         = b.GetSink()
	)

	screeners := []models.Step{
		// Step 1: Confirm bullish candle pattern (hammer, engulfing, piercing, or solid green)
		&steps.BullishCandle{},

		// Step 2: Ensure current volume is at least equal to average volume
		// Higher volume adds conviction to the bullish signal
		&steps.Volume{
			Test: func(currentVolume float64, averageVolume float64) bool {
				return currentVolume >= averageVolume
			},
		},

		// Step 3: Check RSI is in the swing zone (40-60)
		// This range indicates the stock is neither overbought nor oversold,
		// providing room for upward movement while showing some strength
		&steps.Rsi{
			Test: func(rsi []float64) bool {
				var (
					length = len(rsi)
					v      = rsi[length-1]
				)

				return v >= 40.0 && v <= 60.0
			},
		},

		// Step 4: Check if lower Bollinger Band shows support (flat or V-shape)
		// This pattern suggests price is finding support and may bounce higher
		&steps.BollingerBands{
			Test: func(candles []models.Candle, sma []float64, lbb []float64, _ []float64) bool {
				return utils.LowerBollingerBandFlatOrVShape(candles, sma, lbb)
			},
		},
	}

	// Execute all screening steps; if all pass, send stock to output sink
	if steps.Execute(strategyName, stock, screeners) {
		sink <- stock
	}
}
