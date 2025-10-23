package strategy

import (
	"eeye/src/models"
	"eeye/src/steps"
	"eeye/src/utils"
	"math"
)

// BullishMomentumBreakout identifies stocks that are breaking out with strong bullish momentum.
// This is an aggressive strategy that looks for multiple confirmation signals:
//   - Bullish candlestick pattern (strong buying pressure)
//   - RSI breaking above 60 (entering bullish momentum zone from neutral)
//   - Price trading above 50-day EMA (confirming uptrend)
//   - Price breaking above upper Bollinger Band (breakout from volatility envelope)
//   - Properly aligned EMA stack (5 > 13 > 26 > 50 > 200, confirming trend strength)
//
// Trading Logic:
//   - RSI crossing above 60 indicates shift from neutral to bullish momentum
//   - Upper Bollinger Band break suggests expansion and potential for continued move
//   - EMA alignment ensures all timeframes are in bullish agreement
//   - This combination filters for stocks with the highest probability of sustained momentum
//
// Ideal For: Momentum traders looking for explosive breakout opportunities
// Timeframe: Best on daily charts for swing trades, or 4-hour for shorter-term momentum
// Risk Profile: Higher - requires multiple confirmations but targets significant moves
type BullishMomentumBreakout struct {
	models.StrategyBaseImpl
}

// Name returns the strategy identifier.
//
// Returns:
//   - The name of this strategy ("Bullish momentum")
//
//revive:disable-next-line exported
func (b *BullishMomentumBreakout) Name() string {
	return "Bullish momentum"
}

// Execute runs the BullishMomentumBreakout strategy on the given stock.
// It applies five rigorous screening steps:
//  1. BullishCandle: Confirms strong bullish price action
//  2. Rsi: Checks if RSI is breaking above 60 (entering momentum zone)
//  3. Ema: Verifies price is above 50-day EMA (trend confirmation)
//  4. BollingerBands: Confirms breakout above upper band (volatility expansion)
//  5. EmaCrossover: Ensures proper EMA alignment (5>13>26>50>200)
//
// If all five conditions are met, the stock is sent to the strategy's output sink.
//
// Parameters:
//   - stock: The stock to analyze for bullish momentum breakout
//
//revive:disable-next-line exported
func (b *BullishMomentumBreakout) Execute(stock *models.Stock) {
	var (
		strategyName = b.Name()
		sink         = b.GetSink()
	)

	screeners := []models.Step{
		// Step 1: Confirm bullish candlestick pattern
		&steps.BullishCandle{},

		// Step 2: Check if RSI is breaking into momentum zone
		// Previous RSI <= 60 and current RSI > 60 indicates fresh momentum
		&steps.Rsi{
			Test: func(rsi []float64) bool {
				length := len(rsi)

				if length < 2 {
					return false
				}

				// RSI crossing above 60 from below signals momentum shift
				return rsi[length-2] <= 60 &&
					rsi[length-1] > 60
			},
		},

		// Step 3: Verify price is trading above 50-day EMA
		// This confirms the stock is in an uptrend
		&steps.Ema{
			Period: 50,
			Test: func(candles []models.Candle, emas []float64) bool {
				var (
					emaLength    = len(emas)
					candleLength = len(candles)
				)

				// Current close must be above the 50-day EMA
				return candles[candleLength-1].Close > emas[emaLength-1]
			},
		},

		// Step 4: Check for breakout above upper Bollinger Band
		// This indicates price is breaking out of the normal volatility range
		&steps.BollingerBands{
			Test: func(candles []models.Candle, _, _, ubb []float64) bool {
				var (
					ubbLength    = len(ubb)
					candleLength = len(candles)
				)

				// Current high must exceed the upper Bollinger Band
				return candles[candleLength-1].High > ubb[ubbLength-1]
			},
		},

		// Step 5: Verify proper EMA alignment (bullish stack)
		// EMAs should be ordered: 5 > 13 > 26 > 50 > 200
		// This confirms bullish trend across all timeframes
		&steps.EmaCrossover{
			Periods: []int{5, 13, 26, 50, 200},
			Test: func(emas [][]float64) bool {
				// Start with infinity as the "previous" value
				prev := math.Inf(1)

				// Check each EMA is less than the previous (shorter EMAs above longer ones)
				for i := range emas {
					next := utils.Last(emas[i], math.Inf(1))
					if prev < next {
						// If a longer EMA is above a shorter one, alignment is broken
						return false
					}
					prev = next
				}

				// All EMAs are properly aligned (5 > 13 > 26 > 50 > 200)
				return true
			},
		},
	}

	// Execute all screening steps; if all five pass, send stock to output sink
	if steps.Execute(strategyName, stock, screeners) {
		sink <- stock
	}
}
