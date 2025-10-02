// Package models defines the data structures used throughout the trading system.
// It includes types for candlestick data, stock information, and API responses.
package models

import "time"

// Candle represents a single candlestick in a financial chart, containing
// OHLCV (Open, High, Low, Close, Volume) data for a specific time period.
type Candle struct {
	// Symbol identifies the stock this candle belongs to
	Symbol string

	// Open is the opening price for the period
	Open float64

	// Close is the closing price for the period
	Close float64

	// High is the highest price reached during the period
	High float64

	// Low is the lowest price reached during the period
	Low float64

	// Timestamp marks when this candle period started
	Timestamp time.Time

	// Volume is the trading volume during this period
	Volume uint64
}

// RawCandle is a type alias for raw candlestick data received from the API,
// represented as an array of 6 elements containing timestamp, open, high, low,
// close, and volume in that order.
type RawCandle = [6]any

// CandlePayload represents the structure of candle data received in API responses.
// CandlePayload represents the structure of candle data received in API responses.
type CandlePayload struct {
	// Candles is an array of raw candle data
	Candles []RawCandle `json:"candles"`

	// StartTime is the beginning of the data range in ISO format
	StartTime string `json:"start_time"`

	// EndTime is the end of the data range in ISO format
	EndTime string `json:"end_time"`

	// Interval specifies the candle duration in minutes
	Interval uint `json:"interval_in_minutes"`
}

// CandlesResponse represents the complete API response structure for candle data requests.
type CandlesResponse struct {
	// Status indicates the API response status ("success" or "error")
	Status string `json:"status"`

	// Payload contains the actual candle data
	Payload CandlePayload `json:"payload"`
}
