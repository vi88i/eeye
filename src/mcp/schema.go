package mcp

import (
	"encoding/json"
	"log"
	"time"

	"github.com/kaptinlin/jsonschema"
)

//revive:disable-next-line exported
type GetTechnicalDataInput struct {
	Symbol string `json:"symbol"`
}

//revive:disable-next-line exported
type OHLC struct {
	Open  float64 `json:"open"`
	High  float64 `json:"high"`
	Low   float64 `json:"low"`
	Close float64 `json:"close"`
}

//revive:disable-next-line exported
type Indicators struct {
	Volume float64 `json:"volume"`
	EMA5   float64 `json:"ema5"`
	EMA13  float64 `json:"ema13"`
	EMA26  float64 `json:"ema26"`
	EMA50  float64 `json:"ema50"`
	RSI    float64 `json:"rsi"`
}

//revive:disable-next-line exported
type TechnicalData struct {
	Timestamp  time.Time  `json:"date"`
	OHLC       OHLC       `json:"ohlc"`
	Indicators Indicators `json:"indicators"`
}

//revive:disable-next-line exported
type GetTechnicalDataOutput struct {
	Symbol string          `json:"symbol"`
	Data   []TechnicalData `json:"data"`
}

var (
	// GetTechnicalDataInputSchema is the jsonrpc schema for GetTechnicalData tool input
	GetTechnicalDataInputSchema = jsonschema.Object(
		jsonschema.Prop(
			"symbol",
			jsonschema.String(
				jsonschema.MinLen(1),
				jsonschema.Examples("ZOMATO"),
			),
		),
		jsonschema.Required("symbol"),
	)
	// GetTechnicalDataOutputSchema is the jsonrpc schema for GetTechnicalData tool input
	GetTechnicalDataOutputSchema = jsonschema.Object(
		jsonschema.Prop("symbol", jsonschema.String()),
		jsonschema.Prop("data",
			jsonschema.Array(
				jsonschema.Items(
					jsonschema.Object(
						jsonschema.Prop("date",
							jsonschema.String(
								jsonschema.Format("date-time"),
								jsonschema.Description("ISO8601 timestamp, e.g. 2022-10-14T00:00:00"),
							),
						),
						jsonschema.Prop("ohlc",
							jsonschema.Object(
								jsonschema.Description("Open, high, low, and close of the candle"),
								jsonschema.Prop("open", jsonschema.Number()),
								jsonschema.Prop("high", jsonschema.Number()),
								jsonschema.Prop("low", jsonschema.Number()),
								jsonschema.Prop("close", jsonschema.Number()),
							),
						),
						jsonschema.Prop("indicators",
							jsonschema.Object(
								jsonschema.Description("Technical indicators for this timestamp"),
								jsonschema.Prop("volume", jsonschema.Number()),
								jsonschema.Prop("ema5",
									jsonschema.Number(
										jsonschema.Description("Exponential Moving Average with period 5"),
									),
								),
								jsonschema.Prop("ema13",
									jsonschema.Number(
										jsonschema.Description("Exponential Moving Average with period 13"),
									),
								),
								jsonschema.Prop("ema26",
									jsonschema.Number(
										jsonschema.Description("Exponential Moving Average with period 26"),
									),
								),
								jsonschema.Prop("ema50",
									jsonschema.Number(
										jsonschema.Description("Exponential Moving Average with period 50"),
									),
								),
								jsonschema.Prop("rsi",
									jsonschema.Number(
										jsonschema.Description("Relative Strength Index"),
									),
								),
							),
						),
					),
				),
			),
		),
	)
)

// ResolvedSchema stores the schema of tools in JSON format ([]byte)
var ResolvedSchema = map[*jsonschema.Schema][]byte{}

func init() {
	resolve := func(schema *jsonschema.Schema) {
		bytes, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		ResolvedSchema[schema] = bytes
	}

	schemas := []*jsonschema.Schema{
		GetTechnicalDataInputSchema,
		GetTechnicalDataOutputSchema,
	}

	for i := range schemas {
		resolve(schemas[i])
	}
}
