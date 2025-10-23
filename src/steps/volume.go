package steps

import (
	"eeye/src/constants"
	"eeye/src/models"
	"eeye/src/store"
	"eeye/src/utils"
	"log"
)

// Volume screens stocks based on trading volume analysis.
// Compares current trading volume against the average to identify unusual activity
// which often precedes significant price movements.
type Volume struct {
	models.StepBaseImpl
	// Test compares current volume vs average volume to determine screening criteria.
	// Parameters:
	//   - currentVolume: Latest candle's trading volume
	//   - averageVolume: 20-period simple moving average of volume
	// Returns true if the stock passes the volume screening test.
	Test func(currentVolume float64, averageVolume float64) bool
}

//revive:disable-next-line exported
func (v *Volume) Name() string {
	return "Volume screener"
}

//revive:disable-next-line exported
func (v *Volume) Screen(strategy string, stock *models.Stock) bool {
	const (
		Period = 20 // Standard period for volume moving average
	)

	step := v.Name()

	candles, err := store.Get(stock)
	if err != nil {
		return false
	}

	length := len(candles)
	volumeMA := ComputeVolumeMA(candles, Period)
	if length < Period {
		log.Printf("[%v - %v] insufficient candles: %v\n", strategy, step, stock.Symbol)
		return false
	}

	maLength := len(volumeMA)
	return v.TruthyCheck(
		strategy,
		step,
		stock,
		func() bool {
			return v.Test(float64(candles[length-1].Volume), volumeMA[maLength-1])
		},
	)
}

// ComputeVolumeMA calculates the Simple Moving Average of trading volumes.
// Uses a rolling window to compute average volume over the specified period,
// useful for identifying when current volume is significantly above/below normal.
//
// Parameters:
//   - candles: Historical price/volume data
//   - period: Number of candles for average calculation (standard is 20)
//
// Returns:
//   - Slice of volume MA values (empty if insufficient data)
func ComputeVolumeMA(candles []models.Candle, period int) []float64 {
	if len(candles) < period {
		return utils.EmptySlice[float64]()
	}

	// Calculate volume MA using rolling window
	sum := 0.0
	volumeMA := make([]float64, 0, constants.LookBackDays)
	for index := range candles {
		candle := &candles[index]
		sum += float64(candle.Volume)

		// Once we have enough data points, calculate the average
		if index+1 >= period {
			volumeMA = append(volumeMA, sum/float64(period))
			// Remove oldest value from rolling sum to maintain window size
			sum -= float64(candles[index-period+1].Volume)
		}
	}

	return volumeMA
}
