package strategy

import (
	"eeye/src/config"
	"eeye/src/models"
	"eeye/src/steps"
	"log"
	"strings"
	"sync"
)

const (
	strategyName = "Bullish Swing"
)

func worker(in chan *models.Stock, out chan *models.Stock, wg *sync.WaitGroup) {
	for stock := range in {
		if err := steps.Ingestor(stock); err != nil {
			log.Printf("ingestion failed for %v: %v\n", stock.Symbol, err)
			continue
		}

		if err := steps.Extractor(stock); err != nil {
			log.Printf("historical data extraction failed for %v: %v", stock.Symbol, err)
			continue
		}

		screeners := []func() bool{
			steps.VolumeScreener(
				strategyName,
				stock,
				func(currentVolume float64, averageVolume float64) bool {
					return currentVolume >= averageVolume
				},
			),
			steps.RSIScreener(
				strategyName,
				stock,
				func(currentRSI float64) bool {
					return currentRSI >= 40.0 && currentRSI <= 60.0
				},
			),
		}

		if steps.Executor(screeners) {
			out <- stock
		}

		steps.PurgeCache(stock)
	}
	wg.Done()
}

func BullishSwing(stocks []models.Stock) {
	var (
		wg  = sync.WaitGroup{}
		in  = make(chan *models.Stock, config.StepsConfig.Concurrency)
		out = make(chan *models.Stock, config.StepsConfig.Concurrency)
	)

	for i := 1; i <= config.StepsConfig.Concurrency; i++ {
		wg.Add(1)
		go worker(in, out, &wg)
	}

	go func() {
		for _, stock := range stocks {
			in <- &stock
		}
		close(in)
	}()

	go func() {
		defer close(out)
		wg.Wait()
	}()

	filtered := []string{}
	for stock := range out {
		filtered = append(filtered, stock.Symbol)
	}

	if len(filtered) > 0 {
		log.Printf("%v result: \n%v\n", strategyName, strings.Join(filtered, "\n"))
	}
}
