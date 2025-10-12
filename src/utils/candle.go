package utils

import (
	"eeye/src/models"
	"time"
)

// GetTimestamps extracts all the timestamp from the candles array
func GetTimestamps(candles []models.Candle) []time.Time {
	return Map(
		candles,
		func(candle models.Candle) time.Time {
			return candle.Timestamp
		},
	)
}
