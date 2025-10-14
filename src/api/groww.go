package api

import (
	"eeye/src/config"
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

	loc, err := time.LoadLocation(config.DB.Tz)
	if err != nil {
		return empty, fmt.Errorf("unable to load location: %w", err)
	}

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

		timestamp, ok := c[constants.CandleTimestampIndex].(float64)
		if !ok {
			return empty, fmt.Errorf("invalid timestamp format for %s at index %d", stock.Symbol, i)
		}

		open, ok := c[constants.CandleOpenIndex].(float64)
		if !ok {
			return empty, fmt.Errorf("invalid open price format for %s at index %d", stock.Symbol, i)
		}

		high, ok := c[constants.CandleHighIndex].(float64)
		if !ok {
			return empty, fmt.Errorf("invalid high price format for %s at index %d", stock.Symbol, i)
		}

		low, ok := c[constants.CandleLowIndex].(float64)
		if !ok {
			return empty, fmt.Errorf("invalid low price format for %s at index %d", stock.Symbol, i)
		}

		closePrice, ok := c[constants.CandleCloseIndex].(float64)
		if !ok {
			return empty, fmt.Errorf("invalid close price format for %s at index %d", stock.Symbol, i)
		}

		volume, ok := c[constants.CandleVolumeIndex].(float64)
		if !ok {
			return empty, fmt.Errorf("invalid volume format for %s at index %d", stock.Symbol, i)
		}

		ts := time.Unix(int64(timestamp), 0)
		startOfDay := time.Date(
			ts.In(loc).Year(),
			ts.In(loc).Month(),
			ts.In(loc).Day(),
			0, 0, 0, 0,
			loc,
		)

		candle := models.Candle{
			Symbol:    stock.Symbol,
			Timestamp: startOfDay,
			Open:      open,
			High:      high,
			Low:       low,
			Close:     closePrice,
			Volume:    uint64(volume),
		}

		candles = append(candles, candle)
	}

	return candles, nil
}
