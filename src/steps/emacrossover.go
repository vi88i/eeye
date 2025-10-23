package steps

import (
	"eeye/src/models"
	"eeye/src/store"
	"log"
)

// EmaCrossover screens for EMA crossover signals between multiple periods.
// Crossovers occur when a faster EMA crosses above/below a slower EMA,
// indicating potential trend changes.
type EmaCrossover struct {
	models.StepBaseImpl
	// Periods is a list of EMA periods to calculate and compare.
	// Example: [9, 21, 50] for short, medium, and long-term EMAs
	Periods []int
	// Test receives EMA arrays for all periods to check crossover conditions.
	// Parameters:
	//   - emas: 2D slice where emas[i] contains EMA values for Periods[i]
	// Returns true if the crossover condition is met.
	Test func(emas [][]float64) bool
}

//revive:disable-next-line exported
func (e *EmaCrossover) Name() string {
	return "EMA crossover"
}

//revive:disable-next-line exported
func (e *EmaCrossover) Screen(strategy string, stock *models.Stock) bool {
	var emas [][]float64

	step := e.Name()

	candles, err := store.Get(stock)
	if err != nil {
		return false
	}

	// Calculate EMAs for all specified periods
	for i, period := range e.Periods {
		emas = append(emas, ComputeEma(candles, period))
		if len(emas[i]) == 0 {
			log.Printf("[%v - %v] insufficient candles for EMA %v: %v\n", strategy, step, period, stock.Symbol)
			return false
		}
	}

	return e.TruthyCheck(
		strategy,
		step,
		stock,
		func() bool {
			return e.Test(emas)
		},
	)
}
