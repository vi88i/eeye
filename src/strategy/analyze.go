package strategy

import (
	"eeye/src/constants"
	"eeye/src/dataflow"
	"eeye/src/models"
	"eeye/src/store"
	"eeye/src/utils"
	"log"
	"strings"
	"sync"
	"time"
)

// executor processes stocks from the source channel, applies all strategies,
// and sends the results to their respective sinks.
func executor(strategies []models.Strategy, source <-chan *models.Stock) {
	for stock := range source {
		if err := store.Add(stock); err != nil {
			log.Printf("historical data extraction failed for %v: %v\n", stock.Symbol, err)
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
		store.Purge(stock)
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

		// What is the purpose of done?
		// It is mainly used as a signalling channel (or to block code until certain condition)
		// Here we use close(done) to indicate that all strategy workers have shutdown,
		// and we can close all the sinks to initiate aggregator goroutine shutdown.
		//
		// Why we can't return WaitGroup?
		// It is a bad practice to return WaitGroup, use channel for goroutine comms.
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
		wg  = sync.WaitGroup{}
		agg = make(chan *models.StrategyResult, len(strategies))
	)

	for i := range strategies {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res := make([]*models.Stock, 0, constants.StrategyWorkerOutputBufferSize)
			for stock := range strategies[i].GetSink() {
				res = append(res, stock)
			}
			agg <- &models.StrategyResult{Strategy: strategies[i], Stocks: res}
		}()
	}

	go func() {
		// Wait for strategy workers to finish
		<-done

		// Since all strategy workers have stopped, we can safely close all strategy sinks
		// to initiate the shutdown of aggregator goroutines
		for i := range strategies {
			close(strategies[i].GetSink())
		}
	}()

	go func() {
		// Wait for all aggregator routines to shutdown, post which
		// we can close the agg channel
		wg.Wait()

		close(agg)
	}()

	for result := range agg {
		strategyName := result.Strategy.Name()
		symbols := utils.EmptySlice[string]()

		for i := range result.Stocks {
			symbols = append(symbols, result.Stocks[i].Symbol)
		}

		if len(symbols) > 0 {
			log.Printf("%v result: \n%v\n", strategyName, strings.Join(symbols, "\n"))
		} else {
			log.Printf("no stocks satisfy %v\n", strategyName)
		}
	}
}

// Analyze runs all trading strategies on the given list of stocks.
// It executes multiple strategies in parallel and logs their results.
func Analyze() <-chan any {
	done := make(chan any)

	go func() {
		defer close(done)

		start := time.Now()
		stocks := dataflow.GetStocks()

		strategies := []models.Strategy{
			&BullishSwing{},
			&LowerBollingerBandBullish{},
			&EMAFakeBreakdown{period: 50},
			&RSIEntersBullishSwingZone{baseLine: 40, upperBound: 60},
		}

		source, isWorkDone := spawnStrategyWorkers(strategies)
		feeder(stocks, source)
		aggregator(strategies, isWorkDone)
		log.Printf("time taken to complete analysis %s\n", time.Since(start))
	}()

	return done
}
