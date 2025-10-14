package mcp

import (
	"context"
	"eeye/src/db"
	"eeye/src/models"
	"eeye/src/steps"
	"eeye/src/utils"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func getTechnicalData(
	_ context.Context,
	_ *mcp.CallToolRequest,
	input GetTechnicalDataInput,
) (*mcp.CallToolResult, GetTechnicalDataOutput, error) {
	res := GetTechnicalDataInputSchema.Validate(input)
	if !res.IsValid() {
		return nil, GetTechnicalDataOutput{}, fmt.Errorf("schema error: %v", res.Error())
	}

	stock := models.Stock{
		Symbol:   input.Symbol,
		Exchange: "NSE",
		Segment:  "CASH",
		Name:     input.Symbol,
	}
	candles, err := db.FetchAllCandles(&stock)
	if err != nil {
		return nil, GetTechnicalDataOutput{}, fmt.Errorf("db failure: %v", err)
	}

	var (
		timestamps = utils.GetTimestamps(candles)
		ohlc       = utils.Map(
			candles,
			func(candle models.Candle) Ohlc {
				return Ohlc{
					Open:  utils.Round2(candle.Open),
					High:  utils.Round2(candle.High),
					Low:   utils.Round2(candle.Low),
					Close: utils.Round2(candle.Close),
				}
			},
		)
		totalItems = len(timestamps)
		rsi        = utils.PadLeft(steps.ComputeRsi(candles, 14), totalItems, -1)
		ema5       = utils.PadLeft(steps.ComputeEma(candles, 5), totalItems, -1)
		ema13      = utils.PadLeft(steps.ComputeEma(candles, 13), totalItems, -1)
		ema26      = utils.PadLeft(steps.ComputeEma(candles, 26), totalItems, -1)
		ema50      = utils.PadLeft(steps.ComputeEma(candles, 50), totalItems, -1)
		volume     = utils.PadLeft(steps.ComputeVolumeMA(candles, 20), totalItems, -1)
	)

	out := GetTechnicalDataOutput{
		Symbol: input.Symbol,
		Data: func() []TechnicalData {
			data := make([]TechnicalData, 0, totalItems)
			for i := range totalItems {
				data = append(data, TechnicalData{
					Timestamp: timestamps[i],
					Ohlc:      ohlc[i],
					Indicators: Indicators{
						Rsi:    utils.Round2(rsi[i]),
						Ema5:   utils.Round2(ema5[i]),
						Ema13:  utils.Round2(ema13[i]),
						Ema26:  utils.Round2(ema26[i]),
						Ema50:  utils.Round2(ema50[i]),
						Volume: utils.Round2(volume[i]),
					},
				})
			}
			sort.Slice(data, func(i, j int) bool {
				return data[i].Timestamp.After(data[j].Timestamp)
			})

			return data
		}(),
	}

	return nil, out, nil
}

func getOhlcData(
	_ context.Context,
	_ *mcp.CallToolRequest,
	input GetOhlcDataInput,
) (*mcp.CallToolResult, GetOhlcDataOutput, error) {
	res := GetOhlcDataInputSchema.Validate(input)
	if !res.IsValid() {
		return nil, GetOhlcDataOutput{}, fmt.Errorf("schema error: %v", res.Error())
	}

	stock := models.Stock{
		Symbol:   input.Symbol,
		Exchange: "NSE",
		Segment:  "CASH",
		Name:     input.Symbol,
	}
	candles, err := db.FetchAllCandles(&stock)
	if err != nil {
		return nil, GetOhlcDataOutput{}, fmt.Errorf("db failure: %v", err)
	}

	var (
		timestamps = utils.GetTimestamps(candles)
		totalItems = len(timestamps)
	)

	out := GetOhlcDataOutput{
		Symbol: input.Symbol,
		Data: func() []OhlcWithTimestamp {
			data := make([]OhlcWithTimestamp, 0, totalItems)
			for i := range totalItems {
				data = append(data, OhlcWithTimestamp{
					Timestamp: timestamps[i],
					Ohlc: []float64{
						utils.Round2(candles[i].Open),
						utils.Round2(candles[i].High),
						utils.Round2(candles[i].Low),
						utils.Round2(candles[i].Close),
					},
				})
			}
			sort.Slice(data, func(i, j int) bool {
				return data[i].Timestamp.After(data[j].Timestamp)
			})

			return data
		}(),
	}

	return nil, out, nil
}

func addTools(server *mcp.Server) {
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:         "getTechnicalData",
			Title:        "Technical data of symbol",
			Description:  "Gives OHLC, EMA, Volume, RSI etc. indicators",
			InputSchema:  json.RawMessage(ResolvedSchema[GetTechnicalDataInputSchema]),
			OutputSchema: json.RawMessage(ResolvedSchema[GetTechnicalDataOutputSchema]),
		},
		getTechnicalData,
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:         "getOhlcData",
			Title:        "OHLC data of symbol",
			Description:  "Gives OHLC with timestamp for the given symbol",
			InputSchema:  json.RawMessage(ResolvedSchema[GetOhlcDataInputSchema]),
			OutputSchema: json.RawMessage(ResolvedSchema[GetOhlcDataOutputSchema]),
		},
		getOhlcData,
	)
}
