// Package store provides cache service
package store

import (
	"eeye/src/db"
	"eeye/src/models"
	"fmt"
	"log"
	"sync"
)

var cache map[string][]models.Candle
var mu sync.RWMutex

func init() {
	cache = make(map[string][]models.Candle)
}

// Get returns the candles of the given stock
func Get(stock *models.Stock) ([]models.Candle, error) {
	mu.RLock()
	defer mu.RUnlock()

	value, ok := cache[stock.Symbol]
	if !ok {
		return value, fmt.Errorf("unexpected cache miss: %v", stock.Symbol)
	}

	return value, nil
}

// Add retrieves candlestick data for a stock from the database and
// caches it in memory for faster access by other analysis functions. This helps
// prevent repeated database queries for the same data.
func Add(stock *models.Stock) error {
	candles, err := db.FetchAllCandles(stock)
	if err != nil {
		return fmt.Errorf("failed to fetch candles for %v: %w", stock.Symbol, err)
	}

	mu.Lock()
	defer mu.Unlock()
	cache[stock.Symbol] = candles
	return nil
}

// Purge removes the cached candlestick data for a specific stock.
func Purge(stock *models.Stock) {
	mu.Lock()
	defer mu.Unlock()

	_, ok := cache[stock.Symbol]
	if ok {
		log.Printf("purged %v from cache\n", stock.Symbol)
		delete(cache, stock.Symbol)
	}
}
