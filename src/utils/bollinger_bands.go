package utils

import (
	"eeye/src/models"
	"math"
)

// LowerBollingerBandFlatOrVShape detects flat or V-shaped patterns in the lower Bollinger Band.
// These patterns often indicate potential bullish reversal points where the downtrend may be ending.
//
// Pattern Detection:
//   - Flat: Lower band slope is near zero for last 3 periods (stable support)
//   - V-Shape: Lower band slopes down then up (bounce from support)
//
// Parameters:
//   - candles: Historical price data
//   - sma: Simple Moving Average values (middle Bollinger Band)
//   - lbb: Lower Bollinger Band values
//
// Returns:
//   - true if a flat or V-shaped pattern is detected AND price hasn't fully crossed above SMA
//   - false otherwise
//
// Note: Only triggers when price is still near or below SMA to avoid late signals
func LowerBollingerBandFlatOrVShape(candles []models.Candle, sma []float64, lbb []float64) bool {
	const (
		FlatThreshold = 0.0001 // Threshold for detecting flat slope (near zero)
	)

	// Check if price has fully crossed above SMA
	last := len(sma) - 1
	lastCandle := candles[len(candles)-1]
	isBullishCandle := lastCandle.Open >= lastCandle.Close
	hasCandleFullyCrossedSMA := (isBullishCandle && lastCandle.Open > sma[last]) ||
		(lastCandle.Close > sma[last])

	// Only check for pattern if price hasn't fully crossed SMA yet
	if !hasCandleFullyCrossedSMA {
		lbbLen := len(lbb)
		// Set up coordinates for last 3 lower band values
		var (
			x1 = 0.0
			y1 = lbb[lbbLen-3] // Third-to-last value
			x2 = 1.0
			y2 = lbb[lbbLen-2] // Second-to-last value
			x3 = 2.0
			y3 = lbb[lbbLen-1] // Last value
		)

		// Calculate slopes between consecutive points
		slope1 := (y2 - y1) / (x2 - x1) // Slope from point 1 to point 2
		slope2 := (y3 - y2) / (x3 - x2) // Slope from point 2 to point 3

		// Detect patterns
		isFlat := math.Abs(slope1) < FlatThreshold && math.Abs(slope2) < FlatThreshold
		isVShape := slope1 < 0 && slope2 > 0 // First slope down, second slope up
		if isFlat || isVShape {
			return true
		}
	}

	return false
}
