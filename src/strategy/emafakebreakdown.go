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

//nolint:revive
func (e *EMAFakeBreakdown) Name() string {
	return fmt.Sprintf("EMA %v fake breakdown", e.period)
}

//nolint:revive
func (e *EMAFakeBreakdown) Execute(stock *models.Stock) {
	var (
		strategyName = e.Name()
		sink         = e.GetSink()
	)

	screeners := []func() bool{
		steps.BullishCandleScreener(
			strategyName,
			stock,
		),
		steps.EMAFakeBreakdown(
			strategyName,
			stock,
			e.period,
		),
	}

	if steps.Screen(screeners) {
		sink <- stock
	}
}
