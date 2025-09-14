package api

import (
	"eeye/src/constants"
	"eeye/src/models"
	"fmt"
	"log"
	"time"
)

func GetCandles(stock *models.Stock, startTime string, endTime string) ([]models.Candle, error) {
	log.Printf("fetching candles for %v from %v to %v\n", stock.Symbol, startTime, endTime)
	body := models.CandlesResponse{}
	empty := []models.Candle{}
	
	resp, err := Client.
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
		return empty, fmt.Errorf("unauthorized network request")
	}

	if body.Status != "SUCCESS" {
		return empty, fmt.Errorf("internal server error")
	}

	candles := make([]models.Candle, 0, len(body.Payload.Candles))
	for _, rawCandle := range body.Payload.Candles {
		candle := models.Candle{
			Symbol: stock.Symbol,
			Timestamp: time.Unix(int64(rawCandle[constants.CandleTimestampIndex].(float64)), 0),
			Open: rawCandle[constants.CandleOpenIndex].(float64),
			High: rawCandle[constants.CandleHighIndex].(float64),
			Low: rawCandle[constants.CandleLowIndex].(float64),
			Close: rawCandle[constants.CandleCloseIndex].(float64),
			Volume: uint64(rawCandle[constants.CandleVolumeIndex].(float64)),
		}

		candles = append(candles, candle)
	}

	return candles, nil
}
