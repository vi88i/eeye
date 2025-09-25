package steps

import (
	"eeye/src/config"
	"eeye/src/models"
	"eeye/src/utils"
	"fmt"
	"strings"
	"sync"
	"time"
)

// Worker processes a list of stocks in parallel using a configurable number of
// worker goroutines. It applies the given work function to each stock and
// returns a summary of the results.
//
// Parameters:
//   - strategyName: identifier for logging and result reporting
//   - stocks: list of stocks to process
//   - work: function to apply to each stock, receiving input and output channels
//     for stock processing
func Worker(
	strategyName string,
	stocks []models.Stock,
	work func(strategyName string, in chan *models.Stock, out chan *models.Stock),
) string {
	const (
		Delta = 100
	)

	var (
		wg  = sync.WaitGroup{}
		in  = make(chan *models.Stock, config.StepsConfig.Concurrency)
		out = make(chan *models.Stock, config.StepsConfig.Concurrency)
	)

	for i := 1; i <= config.StepsConfig.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			work(strategyName, in, out)
		}()
	}

	go func() {
		var sleep = time.Duration(1000/config.TradingAPIConfig.RateLimit) + Delta

		for i := range stocks {
			in <- &stocks[i]
			time.Sleep(time.Millisecond * sleep)
		}
		close(in)
	}()

	go func() {
		defer close(out)
		wg.Wait()
	}()

	filtered := utils.EmptySlice[string]()
	for stock := range out {
		filtered = append(filtered, stock.Symbol)
	}

	if len(filtered) > 0 {
		return fmt.Sprintf("%v result: \n%v\n", strategyName, strings.Join(filtered, "\n"))
	}

	return fmt.Sprintf("No stocks satisfy %v", strategyName)
}
