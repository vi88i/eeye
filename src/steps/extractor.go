package steps

import (
	"eeye/src/db"
	"eeye/src/models"
	"fmt"
	"log"
	"sync"
)

var cache map[*models.Stock][]models.Candle
var mu sync.RWMutex

func init() {
	cache = make(map[*models.Stock][]models.Candle)
}

func getCachedCandles(stock *models.Stock) ([]models.Candle, error) {
	mu.RLock()
	defer mu.RUnlock()

	value, ok := cache[stock]
	if !ok {
		return value, fmt.Errorf("unexpected cache miss: %v", stock.Symbol)
	}

	return value, nil
}

// Extractor retrieves historical price data for a stock from the database and
// caches it in memory for faster access by other analysis functions. This helps
// prevent repeated database queries for the same data.
func Extractor(stock *models.Stock) error {
	candles, err := db.FetchAllCandles(stock)
	if err != nil {
		return fmt.Errorf("failed to fetch candles for %v: %w", stock.Symbol, err)
	}

	mu.Lock()
	defer mu.Unlock()
	cache[stock] = candles
	return nil
}

// PurgeCache removes the cached candlestick data for a specific stock.
// This is useful when you want to force a refresh of the data from the database.
func PurgeCache(stock *models.Stock) {
	mu.Lock()
	defer mu.Unlock()

	_, ok := cache[stock]
	if ok {
		log.Printf("purged %v\n", stock.Symbol)
		delete(cache, stock)
	}
}
