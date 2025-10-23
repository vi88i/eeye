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

// BollingerBands screens stocks based on Bollinger Band analysis.
// Bollinger Bands consist of a middle band (SMA) and two outer bands (upper and lower)
// that are standard deviations away from the SMA, used to identify overbought/oversold conditions.
type BollingerBands struct {
	models.StepBaseImpl
	// Test is a custom function that receives candles and Bollinger Band data to determine
	// if the stock passes the screening criteria.
	// Parameters:
	//   - candles: Historical price data
	//   - sma: Simple Moving Average values (middle band)
	//   - lbb: Lower Bollinger Band values (SMA - K*stdDev)
	//   - ubb: Upper Bollinger Band values (SMA + K*stdDev)
	// Returns true if the stock passes the screening test.
	Test func(candles []models.Candle, sma []float64, lbb []float64, ubb []float64) bool
}

//revive:disable-next-line exported
func (b *BollingerBands) Name() string {
	return "Bollinger Bands screener"
}

//revive:disable-next-line exported
func (b *BollingerBands) Screen(strategy string, stock *models.Stock) bool {
	const (
		MinPoints = 22 // Minimum candles required (Period + 2 for meaningful analysis)
		Period    = 20 // Standard Bollinger Band period (20-day SMA)
		K         = 2  // Number of standard deviations for band width
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

	// Calculate Bollinger Bands using a rolling window approach
	sum := 0.0
	lbb := make([]float64, 0, len(candles)-Period+1) // Lower Bollinger Band
	ubb := make([]float64, 0, len(candles)-Period+1) // Upper Bollinger Band
	sma := make([]float64, 0, len(candles)-Period+1) // Simple Moving Average (middle band)

	for i := range candles {
		sum += candles[i].Close

		// Once we have enough data points, calculate the bands
		if i+1 >= Period {
			// Calculate SMA (middle band)
			avg := sum / float64(Period)
			sma = append(sma, avg)

			// Calculate variance for standard deviation
			variance := 0.0
			for j := i + 1 - Period; j <= i; j++ {
				diff := candles[j].Close - avg
				variance = variance + diff*diff
			}
			stdDev := math.Sqrt(variance / float64(Period))

			// Calculate lower and upper bands (K standard deviations from SMA)
			lbb = append(lbb, avg-K*stdDev)
			ubb = append(ubb, avg+K*stdDev)

			// Remove oldest value from rolling sum to maintain window size
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
