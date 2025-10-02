package constants

const (
	// HistoricalDataEndpoint is the API endpoint for retrieving historical candle data
	HistoricalDataEndpoint = "/historical/candle/range"
)

// Indices for accessing elements in the raw candle data array returned by the API.
// The API returns an array with elements in the following order:
// [timestamp, open, high, low, close, volume]
const (
	// CandleTimestampIndex is the position of timestamp in the candle array
	CandleTimestampIndex = 0

	// CandleOpenIndex is the position of opening price in the candle array
	CandleOpenIndex = 1

	// CandleHighIndex is the position of highest price in the candle array
	CandleHighIndex = 2

	// CandleLowIndex is the position of lowest price in the candle array
	CandleLowIndex = 3

	// CandleCloseIndex is the position of closing price in the candle array
	CandleCloseIndex = 4

	// CandleVolumeIndex is the position of trading volume in the candle array
	CandleVolumeIndex = 5
)

const (
	// LookBackDays defines the number of days to look back for historical data
	LookBackDays = 1080 // Approximately 3 years of trading days
)

const (
	// MinRequestPerSecond defines the minimum API requests allowed per second
	MinRequestPerSecond = 1

	// MaxRequestPerSecond defines the maximum API requests allowed per second
	MaxRequestPerSecond = 4
)
