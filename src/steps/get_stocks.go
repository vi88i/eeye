package steps

import (
	"eeye/src/api"
	"eeye/src/models"
	"eeye/src/utils"
	"log"
	"strings"
	"time"
)

// Fetch latest available stocks from NSE
func fetchLatestStocksFromNSE() ([]models.Stock, string, error) {
	stocks, lastTradingDay, err := api.DownloadLatestBhavcopy()
	empty := utils.EmptySlice[models.Stock]()
	if err != nil {
		return empty, "", err
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

	Ingestor(stocks, lastTradingDay)
	time.Sleep(time.Second * 2)
	return stocks
}
