package models

// Stock represents a tradable financial instrument with its identifiers
// and market information.
type Stock struct {
	// Symbol is the unique ticker symbol for the stock
	Symbol string

	// Exchange is the market where the stock is traded
	Exchange string

	// Segment represents the market segment (e.g., NSE, BSE)
	Segment string

	// Name is the company's full name
	Name string
}

// NSEStockData represents the structure of stock data fetched from NSE bhavcopy CSV files.
type NSEStockData struct {
	// TckrSymb is the unique ticker symbol for the stock
	Symbol string `csv:"TckrSymb"`

	// SctySrs indicates the type of stock (e.g., EQ, BE etc.)
	Series string `csv:"SctySrs"`

	// FinInstrmNm is the financial instrument name (company's full name)
	Name string `csv:"FinInstrmNm"`

	// ISIN is the International Securities Identification Number
	// Used to differentiate stocks (INE*) from ETFs (INF*)
	ISIN string `csv:"ISIN"`

	// Sgmt is the market segment (e.g., CM for Capital Market)
	Segment string `csv:"Sgmt"`

	// FinInstrmTp is the financial instrument type (e.g., STK for stock)
	InstrumentType string `csv:"FinInstrmTp"`
}
