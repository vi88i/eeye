package strategy

import (
	"eeye/src/models"
	"eeye/src/steps"
	"fmt"
)

// EMAFakeBreakdown strategy
type EMAFakeBreakdown struct {
	models.StrategyBaseImpl
	period int
}

//revive:disable-next-line exported
func (e *EMAFakeBreakdown) Name() string {
	return fmt.Sprintf("EMA %v fake breakdown", e.period)
}

//revive:disable-next-line exported
func (e *EMAFakeBreakdown) Execute(stock *models.Stock) {
	var (
		strategyName = e.Name()
		sink         = e.GetSink()
	)

	screeners := []models.Step{
		&steps.BullishCandle{},
		&steps.EMA{
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
