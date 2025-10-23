package strategy

import (
	"eeye/src/models"
	"eeye/src/steps"
	"eeye/src/utils"
)

// LowerBollingerBandBullish identifies stocks showing bullish reversal patterns based on
// the lower Bollinger Band's shape. This is a focused strategy that looks for:
//   - Bullish candlestick patterns indicating potential reversal
//   - Lower Bollinger Band forming flat or V-shaped patterns (support zones)
//
// Pattern Interpretation:
//   - Flat Lower Band: Indicates price has found a stable support level
//   - V-Shaped Lower Band: Suggests a bounce from support is occurring
//
// These patterns often precede bullish price movements as they indicate downward momentum
// is slowing or reversing at the lower band support level.
//
// Ideal For: Short-term traders looking for reversal opportunities at Bollinger Band support
// Timeframe: Works on all timeframes, commonly used on daily charts
// Risk Profile: Low to Medium - simple pattern with clear support level
type LowerBollingerBandBullish struct {
	models.StrategyBaseImpl
}

// Name returns the strategy identifier.
//
// Returns:
//   - The name of this strategy ("Lower Bollinger Band Bullish")
//
//revive:disable-next-line exported
func (l *LowerBollingerBandBullish) Name() string {
	return "Lower Bollinger Band Bullish"
}

// Execute runs the LowerBollingerBandBullish strategy on the given stock.
// It applies two screening steps:
//  1. BullishCandle: Confirms a bullish price reversal pattern
//  2. BollingerBands: Checks if lower band shows flat or V-shape (support formation)
//
// If both conditions are met, the stock is sent to the strategy's output sink.
//
// Parameters:
//   - stock: The stock to analyze for lower Bollinger Band bullish pattern
//
//revive:disable-next-line exported
func (l *LowerBollingerBandBullish) Execute(stock *models.Stock) {
	var (
		strategyName = l.Name()
		sink         = l.GetSink()
	)

	screeners := []models.Step{
		// Step 1: Verify presence of bullish candlestick pattern
		&steps.BullishCandle{},

		// Step 2: Check if lower Bollinger Band is forming a bullish pattern
		// Flat pattern: lower band has stabilized (support holding)
		// V-shape pattern: lower band turned upward (bounce in progress)
		&steps.BollingerBands{
			Test: func(candles []models.Candle, sma []float64, lbb []float64, _ []float64) bool {
				return utils.LowerBollingerBandFlatOrVShape(candles, sma, lbb)
			},
		},
	}

	// Execute all screening steps; if both pass, send stock to output sink
	if steps.Execute(strategyName, stock, screeners) {
		sink <- stock
	}
}
