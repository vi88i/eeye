package strategy

import (
	"eeye/src/models"
	"eeye/src/steps"
	"fmt"
)

// EmaFakeBreakdown strategy
type EmaFakeBreakdown struct {
	models.StrategyBaseImpl
	period int
}

//revive:disable-next-line exported
func (e *EmaFakeBreakdown) Name() string {
	return fmt.Sprintf("EMA %v fake breakdown", e.period)
}

//revive:disable-next-line exported
func (e *EmaFakeBreakdown) Execute(stock *models.Stock) {
	var (
		strategyName = e.Name()
		sink         = e.GetSink()
	)

	screeners := []models.Step{
		&steps.BullishCandle{},
		&steps.Ema{
			Period: e.period,
			Test: func(candles []models.Candle, ema []float64) bool {
				candleLength := len(candles)
				emaLength := len(ema)
				return (candles[candleLength-1].Low <= ema[emaLength-1]) &&
					(candles[candleLength-1].High > ema[emaLength-1])
			},
		},
	}

	if steps.Execute(strategyName, stock, screeners) {
		sink <- stock
	}
}
