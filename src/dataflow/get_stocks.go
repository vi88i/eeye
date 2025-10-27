// Package dataflow helps in fetching stocks data from NSE
// and fetching latest candles for each stock from Groww
package dataflow

import (
	"eeye/src/api"
	"eeye/src/models"
	"eeye/src/utils"
	"log"
	"strings"
)

// fetchLatestStocksFromNSE fetches latest available stocks from NSE
func fetchLatestStocksFromNSE() ([]models.Stock, string, error) {
	log.Printf("fetching data from NSE")
	stocks, lastTradingDay, err := api.DownloadLatestBhavcopy()
	empty := utils.EmptySlice[models.Stock]()
	if err != nil {
		return empty, "", err
	}

	filtered := make([]models.Stock, 0, len(stocks))
	for _, s := range stocks {
		// Filter for equity stocks in capital market segment
		// Exclude ETFs by checking ISIN prefix (stocks start with INE, ETFs with INF)
		if strings.TrimSpace(s.Segment) == "CM" &&
			strings.TrimSpace(s.InstrumentType) == "STK" &&
			strings.TrimSpace(s.Series) == "EQ" &&
			strings.HasPrefix(strings.TrimSpace(s.ISIN), "INE") {
			filtered = append(filtered, models.Stock{
				Symbol:   strings.TrimSpace(s.Symbol),
				Name:     strings.TrimSpace(s.Name),
				Exchange: "NSE",
				Segment:  "CASH",
			})
		}
	}

	log.Printf("fetched %d stocks from NSE\n", len(filtered))
	return filtered, lastTradingDay, nil
}

// GetStocks retrieves the list of available stocks from an external source or
// fallbacks to distinct stock symbols in the database
func GetStocks() []models.Stock {
	stocks, lastTradingDay, err := fetchLatestStocksFromNSE()
	if err != nil {
		log.Fatal(err)
	}

	ingestor(stocks, lastTradingDay)
	return stocks
}
