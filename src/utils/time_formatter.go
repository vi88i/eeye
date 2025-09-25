package utils

import (
	"eeye/src/constants"
	"time"
)

// GetFormattedTimestamp formats a time.Time value according to the application's
// standard timestamp format defined in constants.TimestampFmt.
func GetFormattedTimestamp(t time.Time) string {
	return t.Format(constants.TimestampFmt)
}
