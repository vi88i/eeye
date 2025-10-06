package steps

import (
	"eeye/src/models"
	"eeye/src/store"
	"log"
)

// EMACrossover helps to check if the EMAs of different periods are in crossover state
type EMACrossover struct {
	Periods []int
	Test    func(emas [][]float64) bool
}

//revive:disable-next-line exported
func (e *EMACrossover) Name() string {
	return "EMA crossover"
}

//revive:disable-next-line exported
func (e *EMACrossover) Screen(strategy string, stock *models.Stock) bool {
	var emas [][]float64

	step := e.Name()

	candles, err := store.Get(stock)
	if err != nil {
		return false
	}

	for i, period := range e.Periods {
		emas = append(emas, computeEMA(candles, period))
		if len(emas[i]) == 0 {
			log.Printf("[%v - %v] insufficient candles for EMA %v: %v\n", strategy, step, period, stock.Symbol)
			return false
		}
	}

	test := e.Test(emas)
	if !test {
		log.Printf("[%v - %v] test failed: %v\n", strategy, step, stock.Symbol)
	}
	return test
}
