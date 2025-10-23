package utils

import (
	"eeye/src/constants"
	"time"
)

// GetFormattedTimestamp formats a time.Time value according to the application's standard format.
// Uses the timestamp format defined in constants.TimestampFmt for consistency across the application.
//
// Parameters:
//   - t: The time value to format
//
// Returns:
//   - String representation of the timestamp in the standard application format
//
// Example:
//
//	formatted := GetFormattedTimestamp(time.Now())
//	// Returns timestamp formatted according to constants.TimestampFmt
func GetFormattedTimestamp(t time.Time) string {
	return t.Format(constants.TimestampFmt)
}
