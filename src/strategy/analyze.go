package strategy

import (
	"eeye/src/constants"
	"eeye/src/models"
	"eeye/src/steps"
	"eeye/src/utils"
	"log"
	"strings"
	"sync"
)

// executor processes stocks from the source channel, applies all strategies,
// and sends the results to their respective sinks.
func executor(strategies []models.Strategy, source <-chan *models.Stock) {
	for stock := range source {
		if err := steps.Extractor(stock); err != nil {
			log.Printf("historical data extraction failed for %v: %v", stock.Symbol, err)
			continue
		}

		wg := sync.WaitGroup{}
		for i := range strategies {
			wg.Add(1)
			go func() {
				defer wg.Done()
				strategies[i].Execute(stock)
			}()
		}

		wg.Wait()
		steps.PurgeCache(stock)
	}
}

// spawnStrategyWorkers initializes worker goroutines to process stocks using the provided strategies.
// done channel is used to signal when all workers have completed processing.
func spawnStrategyWorkers(strategies []models.Strategy) (chan *models.Stock, chan any) {
	var (
		source = make(chan *models.Stock, constants.StrategyWorkerInputBufferSize)
		done   = make(chan any)
	)

	go func() {
		wg := sync.WaitGroup{}

		for range constants.NumOfStrategyWorkers {
			wg.Add(1)
			go func() {
				defer wg.Done()
				executor(strategies, source)
			}()
		}

		// ideally it is not a good practice to send WaitGroup outside of the function
		// since the done channel is only used for signaling, we can safely close it here
		// without worrying about sending values to it
		wg.Wait()
		close(done)
	}()

	return source, done
}

// feeder sends stocks to the source channel for processing and closes the channel when done.
func feeder(stocks []models.Stock, source chan<- *models.Stock) {
	go func() {
		for i := range stocks {
			source <- &stocks[i]
		}
		close(source)
	}()
}

// aggregator collects results from all strategies and logs them once processing is complete.
func aggregator(strategies []models.Strategy, done <-chan any) {
	var (
		agg = make(map[models.Strategy][]*models.Stock)
		wg  = sync.WaitGroup{}
	)

	for i := range strategies {
		agg[strategies[i]] = make([]*models.Stock, 0, constants.StrategyWorkerOutputBufferSize)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for stock := range strategies[i].GetSink() {
				agg[strategies[i]] = append(agg[strategies[i]], stock)
			}
		}()
	}

	// Wait for strategy workers to finish
	<-done

	// Since all workers have stopped, we can safely close all strategy sinks
	// to initiate the shutdown of aggregator goroutines
	for i := range strategies {
		close(strategies[i].GetSink())
	}

	// Even though all strategy workers have finished,
	// there could still be data in the sinks that need to be read (we've used buffered channels for sinks).
	wg.Wait()

	log.Println("================= Strategy Results =================")
	for strategy, stocks := range agg {
		strategyName := strategy.Name()
		symbols := utils.EmptySlice[string]()

		for i := range stocks {
			symbols = append(symbols, stocks[i].Symbol)
		}

		if len(symbols) > 0 {
			log.Printf("%v result: \n%v\n", strategyName, strings.Join(symbols, "\n"))
		} else {
			log.Printf("No stocks satisfy %v", strategyName)
		}
	}
}

// Analyze runs all trading strategies on the given list of stocks.
// It executes multiple strategies in parallel and logs their results.
func Analyze() chan any {
	isAnalysisDone := make(chan any)

	go func() {
		stocks := steps.GetStocks()

		strategies := []models.Strategy{
			&BullishSwing{},
			&LowerBollingerBandBullish{},
			&EMAFakeBreakdown{period: 50},
			&RSIEntersBullishSwingZone{baseLine: 40, upperBound: 60},
		}

		source, done := spawnStrategyWorkers(strategies)
		feeder(stocks, source)
		aggregator(strategies, done)

		close(isAnalysisDone)
	}()

	return isAnalysisDone
}
