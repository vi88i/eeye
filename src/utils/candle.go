package utils

import (
	"eeye/src/models"
	"time"
)

// GetTimestamps extracts all timestamps from a slice of candles.
// This is useful for time-series analysis and charting operations.
//
// Parameters:
//   - candles: Slice of candle data to extract timestamps from
//
// Returns:
//   - Slice of time.Time values representing candle timestamps
//
// Example:
//
//	timestamps := GetTimestamps(candles)
//	// Can be used for plotting time-series data
func GetTimestamps(candles []models.Candle) []time.Time {
	return Map(
		candles,
		func(candle models.Candle) time.Time {
			return candle.Timestamp
		},
	)
}
