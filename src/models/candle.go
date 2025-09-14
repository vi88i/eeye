package models

import "time"

type Candle struct {
	Symbol    string
	Open      float64
	Close     float64
	High      float64
	Low       float64
	Timestamp time.Time
	Volume    uint64
}

type RawCandle = [6]any

type CandlePayload struct {
	Candles   []RawCandle `json:"candles"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Interval  uint `json:"interval_in_minutes"`
}

type CandlesResponse struct {
	Status  string `json:"status"`
	Payload CandlePayload `json:"payload"`
}
