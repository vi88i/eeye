package constants

const (
	HistoricalDataEndpoint = "/historical/candle/range"
)

// Position of attr in the API response
const (
	CandleTimestampIndex = 0
	CandleOpenIndex = 1
	CandleHighIndex = 2
	CandleLowIndex = 3
	CandleCloseIndex = 4
	CandleVolumeIndex = 5
)

const (
	LookBackDays = 1080
)

const (
	MinRequestPerSecond = 1
	MaxRequestPerSecond = 10
)
