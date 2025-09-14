package utils

import (
	"eeye/src/constants"
	"time"
)

func GetFormattedTimestamp(t time.Time) string {
	return t.Format(constants.TimestampFmt)
}
