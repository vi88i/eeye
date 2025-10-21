// Package steps implements various technical analysis indicators and trading strategies.
// It provides functions for analyzing price patterns, calculating technical indicators,
// and screening stocks based on specific criteria.
package steps

import (
	"eeye/src/models"
	"eeye/src/store"
	"log"
	"math"
)

// BollingerBands creates a function to screen stocks based on their bollinger band values
type BollingerBands struct {
	models.StepBaseImpl
	Test func(candles []models.Candle, sma []float64, lbb []float64, ubb []float64) bool
}

//revive:disable-next-line exported
func (b *BollingerBands) Name() string {
	return "Bollinger Bands screener"
}

//revive:disable-next-line exported
func (b *BollingerBands) Screen(strategy string, stock *models.Stock) bool {
	const (
		MinPoints = 22
		Period    = 20
		K         = 2
	)

	step := b.Name()

	candles, err := store.Get(stock)
	if err != nil {
		return false
	}

	length := len(candles)
	if length < MinPoints {
		log.Printf("[%v - %v] insufficient candles: %v\n", strategy, step, stock.Symbol)
		return false
	}

	sum := 0.0
	lbb := make([]float64, 0, len(candles)-Period+1)
	ubb := make([]float64, 0, len(candles)-Period+1)
	sma := make([]float64, 0, len(candles)-Period+1)
	for i := range candles {
		sum += candles[i].Close
		if i+1 >= Period {
			avg := sum / float64(Period)
			sma = append(sma, avg)
			variance := 0.0
			for j := i + 1 - Period; j <= i; j++ {
				diff := candles[j].Close - avg
				variance = variance + diff*diff
			}
			stdDev := math.Sqrt(variance / float64(Period))
			lbb = append(lbb, avg-K*stdDev)
			ubb = append(ubb, avg+K*stdDev)
			sum -= candles[i+1-Period].Close
		}
	}

	return b.TruthyCheck(
		strategy,
		step,
		stock,
		func() bool {
			return b.Test(candles, sma, lbb, ubb)
		},
	)
}
