package steps

import (
	"eeye/src/models"
	"sync"
)

// Execute runs multiple screening steps concurrently and returns true only if all steps pass.
// This allows combining multiple technical analysis conditions that must all be satisfied.
//
// Parameters:
//   - strategy: Name of the trading strategy being executed
//   - stock: Stock to be screened
//   - screeners: Ordered list of Step implementations to execute
//
// Returns:
//   - true if ALL screeners return true (AND logic)
//   - false if ANY screener returns false
//
// Note: Steps are executed concurrently for performance, but the result requires all to pass.
func Execute(strategy string, stock *models.Stock, screeners []models.Step) bool {
	var (
		wg  = sync.WaitGroup{}
		out = make(chan bool)
	)

	// Execute all screeners concurrently
	for i := range screeners {
		wg.Add(1)

		go func() {
			defer wg.Done()
			v := screeners[i].Screen(strategy, stock)
			out <- v
		}()
	}

	// Close channel when all screeners complete
	go func() {
		defer close(out)
		wg.Wait()
	}()

	// Aggregate results with AND logic (all must be true)
	res := true
	for v := range out {
		res = res && v
	}

	return res
}
