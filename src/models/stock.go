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
	// Symbol is the unique ticker symbol for the stock
	Symbol string `csv:"Symbol"`

	// Indicates the type of stock (e.g., EQ, BE etc.)
	Series string `csv:"Series"`

	// Name is the company's full name
	Name string `csv:"Security Name"`

	// Listed | Permitted, Listed means formally listed on NSE
	// Permitted means special grant given (avoid such stocks)
	Category string `csv:"Category"`
}
