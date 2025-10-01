package api

import (
	"eeye/src/constants"
	"eeye/src/models"
	"eeye/src/utils"
	"fmt"
	"log"
	"time"
)

// GetCandles retrieves candlestick data for a given stock within a specified time range.
// It returns an array of Candle objects containing OHLCV data. If startTime equals endTime,
// or if there's an error in fetching data, it returns an empty slice and the error if any.
func GetCandles(stock *models.Stock, startTime string, endTime string) ([]models.Candle, error) {
	log.Printf("fetching candles for %v from %v to %v\n", stock.Symbol, startTime, endTime)
	var (
		body  = models.CandlesResponse{}
		empty = utils.EmptySlice[models.Candle]()
	)

	if startTime >= endTime {
		log.Printf("start time and end time are the same for %v, returning empty slice\n", stock.Symbol)
		return empty, nil
	}

	resp, err := GrowwClient.
		R().
		SetQueryParam("exchange", stock.Exchange).
		SetQueryParam("segment", stock.Segment).
		SetQueryParam("trading_symbol", stock.Symbol).
		SetQueryParam("start_time", startTime).
		SetQueryParam("end_time", endTime).
		SetQueryParam("interval_in_minutes", "1440").
		SetResult(&body).
		Get(constants.HistoricalDataEndpoint)

	if err != nil {
		return empty, fmt.Errorf("failed to get candles: %w", err)
	}

	if resp.IsError() {
		return empty, fmt.Errorf("unauthorized network request: %v", resp.StatusCode())
	}

	if body.Status != "SUCCESS" {
		return empty, fmt.Errorf("internal server error")
	}

	candles := make([]models.Candle, 0, len(body.Payload.Candles))
	for i := range body.Payload.Candles {
		c := &body.Payload.Candles[i]
		candle := models.Candle{
			Symbol:    stock.Symbol,
			Timestamp: time.Unix(int64(c[constants.CandleTimestampIndex].(float64)), 0),
			Open:      c[constants.CandleOpenIndex].(float64),
			High:      c[constants.CandleHighIndex].(float64),
			Low:       c[constants.CandleLowIndex].(float64),
			Close:     c[constants.CandleCloseIndex].(float64),
			Volume:    uint64(c[constants.CandleVolumeIndex].(float64)),
		}

		candles = append(candles, candle)
	}

	return candles, nil
}
