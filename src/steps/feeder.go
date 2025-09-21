package steps

import (
	"eeye/src/config"
	"eeye/src/models"
	"time"
)

const (
	Delta = 100 // additional sleep time to avoid 429 errors
)

func Feeder(in chan *models.Stock, stocks []models.Stock) {
	go func() {
		var sleep = time.Duration(1000 / config.TradingAPIConfig.RateLimit) + Delta
		for i := range stocks {
			in <- &stocks[i]
			time.Sleep(time.Millisecond * sleep)
		}
		close(in)
	}()
}
