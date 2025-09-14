package strategy

import (
	"eeye/src/config"
	"eeye/src/models"
	"eeye/src/steps"
	"log"
	"sync"
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

		ok := steps.VolumeScreener(
			stock,
			func(currentVolume float64, averageVolume float64) bool {
				return currentVolume >= averageVolume 
			},
		)

		if ok {
			out <- stock
		}
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

	for _, stock := range stocks {
		in <- &stock
	}
	close(in)

	go func() {
		wg.Wait()
		close(out)
	}()

	for stock := range out {
		log.Printf("Received %v\n", stock.Symbol)
	}
}
