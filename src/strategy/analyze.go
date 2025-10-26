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

	progressbar "github.com/schollz/progressbar/v3"
)

// executor processes stocks from the source channel, applies all strategies concurrently,
// and sends the results to their respective sinks.
//
// For each stock received from the source channel:
//  1. Fetch and cache historical data in the store
//  2. Execute all strategies concurrently (each in its own goroutine)
//  3. Wait for all strategies to complete
//  4. Clean up the stock data from the store
//
// This design allows multiple strategies to analyze the same stock simultaneously,
// maximizing throughput while ensuring proper cleanup after analysis.
//
// Parameters:
//   - strategies: List of strategies to apply to each stock
//   - source: Channel providing stocks to analyze
func executor(strategies []models.Strategy, source <-chan *models.Stock, bar *progressbar.ProgressBar) {
	// Process each stock from the source channel until it's closed
	for stock := range source {
		// Fetch and cache historical data for this stock
		if err := store.Add(stock); err != nil {
			log.Printf("historical data extraction failed for %v: %v\n", stock.Symbol, err)
			_ = bar.Add(1)
			continue
		}

		// Execute all strategies concurrently for this stock
		wg := sync.WaitGroup{}
		for i := range strategies {
			wg.Go(func() {
				// Each strategy runs independently on the same stock data
				strategies[i].Execute(stock)
			})
		}

		// Wait for all strategies to finish analyzing this stock
		wg.Wait()

		// Clean up cached data for this stock to free memory
		store.Purge(stock)
		_ = bar.Add(1)
	}
}

// spawnStrategyWorkers initializes a pool of worker goroutines to process stocks concurrently.
// This function creates a worker pool pattern where multiple workers pull stocks from the same
// source channel and execute all strategies on each stock.
//
// Architecture:
//   - Creates a buffered input channel for stocks (source)
//   - Spawns N worker goroutines (defined by constants.NumOfStrategyWorkers)
//   - Each worker runs the executor function, pulling stocks from the source channel
//   - Returns a done channel that signals when all workers have finished
//
// The done channel pattern:
//   - Used instead of returning a WaitGroup (which is an anti-pattern)
//   - Allows other goroutines to wait for all workers to complete
//   - Closed when all workers have finished processing
//   - Enables graceful shutdown of downstream components (aggregators)
//
// Parameters:
//   - strategies: List of strategies that each worker will execute on stocks
//
// Returns:
//   - source: Buffered channel to send stocks for processing
//   - done: Signal channel that closes when all workers have finished
func spawnStrategyWorkers(strategies []models.Strategy, numOfStocks int) (chan *models.Stock, chan any) {
	var (
		// Buffered channel to prevent blocking when sending stocks
		source = make(chan *models.Stock, constants.StrategyWorkerInputBufferSize)
		// Signal channel for completion notification
		done = make(chan any)
	)

	go func() {
		wg := sync.WaitGroup{}
		bar := utils.GetProgressTracker(numOfStocks, "Analyzing stocks...")

		// Spawn N worker goroutines to process stocks in parallel
		for range constants.NumOfStrategyWorkers {
			wg.Go(func() {
				// Each worker runs the executor, pulling from the shared source channel
				executor(strategies, source, bar)
			})
		}

		// Wait for all workers to finish processing
		// This happens after the source channel is closed and all stocks are processed
		wg.Wait()

		// Signal completion by closing the done channel
		// This allows downstream aggregators to begin their shutdown sequence
		// Note: We use a channel instead of returning WaitGroup because:
		// - Channels are the idiomatic way for goroutine communication in Go
		// - Exposing WaitGroup creates tight coupling and is considered bad practice
		// - Channels provide a cleaner API for signaling completion
		close(done)
	}()

	return source, done
}

// feeder sends stocks to the source channel for processing and closes the channel when done.
// This function acts as a producer in the producer-consumer pattern, feeding stocks to the
// worker pool for parallel analysis.
//
// Process:
//  1. Iterates through all stocks in the provided slice
//  2. Sends each stock to the source channel (consumed by workers)
//  3. Closes the source channel when all stocks have been sent
//
// Closing the channel signals to workers that no more stocks will be sent,
// allowing them to finish processing and exit gracefully.
//
// Parameters:
//   - stocks: Slice of stocks to feed to the worker pool
//   - source: Send-only channel where stocks are sent for processing
func feeder(stocks []models.Stock, source chan<- *models.Stock) {
	go func() {
		// Send each stock to the worker pool
		for i := range stocks {
			source <- &stocks[i]
		}

		// Close the channel to signal no more stocks will be sent
		// This allows workers to exit their range loops and shutdown gracefully
		close(source)
	}()
}

// aggregator collects results from all strategies and logs them once processing is complete.
// This function implements a fan-in pattern, collecting results from multiple strategy sinks
// into a single aggregation point for reporting.
//
// Architecture:
//  1. For each strategy, spawn a goroutine to collect stocks from its sink
//  2. Wait for all strategy workers to finish (via done channel)
//  3. Close all strategy sinks to signal aggregators to finish
//  4. Collect all strategy results and log them
//
// Shutdown Sequence:
//   - done channel closes → all workers finished processing
//   - Strategy sinks close → aggregators finish collecting
//   - Aggregation channel closes → final results are logged
//
// This coordinated shutdown ensures all results are collected before logging.
//
// Parameters:
//   - strategies: List of strategies whose results need to be collected
//   - done: Signal channel indicating when strategy workers have finished
func aggregator(strategies []models.Strategy, done <-chan any) {
	var (
		wg  = sync.WaitGroup{}
		agg = make(chan *models.StrategyResult, len(strategies))
	)

	// Spawn a goroutine for each strategy to collect its results
	for i := range strategies {
		wg.Go(func() {
			// Collect all stocks that passed this strategy's screening
			res := make([]*models.Stock, 0, constants.StrategyWorkerOutputBufferSize)
			for stock := range strategies[i].GetSink() {
				res = append(res, stock)
			}

			// Send the complete result set to the aggregation channel
			agg <- &models.StrategyResult{Strategy: strategies[i], Stocks: res}
		})
	}

	// Shutdown coordinator goroutine
	go func() {
		// Wait for all strategy workers to finish processing stocks
		<-done

		// Since all strategy workers have stopped, we can safely close all strategy sinks
		// This signals to the aggregator goroutines (above) that no more stocks will arrive
		for i := range strategies {
			close(strategies[i].GetSink())
		}
	}()

	// Aggregation channel closer goroutine
	go func() {
		// Wait for all aggregator goroutines to finish collecting results
		wg.Wait()

		// Close the aggregation channel to signal that all results are in
		// This allows the main loop below to exit
		close(agg)
	}()

	// Collect and log results from all strategies
	for result := range agg {
		strategyName := result.Strategy.Name()
		symbols := utils.EmptySlice[string]()

		// Extract symbols from stocks that passed this strategy
		for i := range result.Stocks {
			symbols = append(symbols, result.Stocks[i].Symbol)
		}

		// Log results
		if len(symbols) > 0 {
			log.Printf("%v result: \n%v\n", strategyName, strings.Join(symbols, "\n"))
		} else {
			log.Printf("no stocks satisfy %v\n", strategyName)
		}
	}
}

// Analyze orchestrates the execution of all trading strategies on the stock universe.
// This is the main entry point for strategy analysis, coordinating the entire pipeline:
//  1. Fetch all stocks from the data source
//  2. Initialize and configure all trading strategies
//  3. Spawn worker pool to process stocks concurrently
//  4. Feed stocks to the worker pool
//  5. Aggregate and log results from all strategies
//  6. Report total execution time
//
// Concurrency Model:
//   - Multiple worker goroutines process stocks in parallel
//   - Each stock is analyzed by all strategies concurrently
//   - Results are collected via a fan-in aggregation pattern
//
// Active Strategies:
//   - BullishSwing: RSI 40-60 with Bollinger Band support
//   - LowerBollingerBandBullish: Flat or V-shaped lower band patterns
//   - EmaFakeBreakdown: Fake breakdown below 50-day EMA
//   - FakeBreakdown: Fake breakdown below support levels (5-day window, 1% tolerance)
//   - RsiEntersBullishSwingZone: RSI crossing into 40-60 range
//   - BullishMomentumBreakout: Strong momentum with EMA alignment
//
// Returns:
//   - Signal channel that closes when analysis is complete
//     This allows callers to wait for completion if needed
func Analyze() <-chan any {
	done := make(chan any)

	go func() {
		defer close(done)

		start := time.Now()

		// Fetch all stocks from the data source
		stocks := dataflow.GetStocks()

		// Initialize all trading strategies with their configurations
		strategies := []models.Strategy{
			// Swing trading strategy looking for balanced momentum
			&BullishSwing{},

			// Simple Bollinger Band reversal strategy
			&LowerBollingerBandBullish{},

			// Fake breakdown at 50-day EMA (dynamic support)
			&EmaFakeBreakdown{period: 50},

			// Fake breakdown at static support levels
			// Window: 5 periods for recent support identification
			// Tolerance: 1% price clustering for level formation
			// Strength: 3 minimum touches to confirm level validity
			&FakeBreakdown{
				Window:    5,
				Tolerance: 0.01,
				Strength:  3,
			},

			// RSI momentum shift detection
			// baseLine: 40 (minimum RSI to enter swing zone)
			// upperBound: 60 (maximum RSI to avoid overbought)
			&RsiEntersBullishSwingZone{baseLine: 40, upperBound: 60},

			// Aggressive breakout strategy with multiple confirmations
			&BullishMomentumBreakout{},
		}

		// Set up concurrent processing pipeline
		source, isWorkDone := spawnStrategyWorkers(strategies, len(stocks))
		feeder(stocks, source)
		aggregator(strategies, isWorkDone)

		// Log performance metrics
		log.Printf("time taken to complete analysis %s\n", time.Since(start))
	}()

	return done
}
