package steps

import (
	"eeye/src/api"
	"eeye/src/db"
	"eeye/src/models"
	"eeye/src/utils"
	"log"
	"strings"
)

// Fetch latest available stocks from NSE
func fetchLatestStocksFromNSE() ([]models.Stock, error) {
	stocks, err := api.DownloadLatestBhavcopy()
	empty := utils.EmptySlice[models.Stock]()
	if err != nil {
		return empty, err
	}

	filtered := make([]models.Stock, 0, len(stocks))
	for _, s := range stocks {
		if strings.TrimSpace(s.Series) == "EQ" && strings.TrimSpace(s.Category) == "Listed" {
			filtered = append(filtered, models.Stock{
				Symbol:   strings.TrimSpace(s.Symbol),
				Name:     strings.TrimSpace(s.Name),
				Exchange: "NSE",
				Segment:  "CASH",
			})
		}
	}

	log.Printf("Fetched %d stocks from NSE\n", len(filtered))
	return filtered, nil
}

// Fetch distinct stock symbols from the DB
func fetchDistinctStocksFromDB() ([]models.Stock, error) {
	stocks, err := db.FetchAllStocks()
	empty := utils.EmptySlice[models.Stock]()
	if err != nil {
		return empty, err
	}

	log.Printf("Fetched %d stocks from DB\n", len(stocks))
	return stocks, nil
}

// GetStocks retrieves the list of available stocks from an external source or
// fallbacks to distinct stock symbols in the database
func GetStocks() []models.Stock {
	if stocks, err := fetchLatestStocksFromNSE(); err == nil {
		return stocks
	}

	if stocks, err := fetchDistinctStocksFromDB(); err == nil {
		return stocks
	}

	log.Fatal("Failed to fetch stocks")
	return utils.EmptySlice[models.Stock]()
}
