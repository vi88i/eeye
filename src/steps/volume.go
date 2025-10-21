package steps

import (
	"eeye/src/constants"
	"eeye/src/models"
	"eeye/src/store"
	"eeye/src/utils"
	"log"
)

// Volume creates a function that screens stocks based on their trading volume.
// It compares the current volume against the average volume using a custom screening
// function to identify significant volume patterns.
type Volume struct {
	models.StepBaseImpl
	Test func(currentVolume float64, averageVolume float64) bool
}

//revive:disable-next-line exported
func (v *Volume) Name() string {
	return "Volume screener"
}

//revive:disable-next-line exported
func (v *Volume) Screen(strategy string, stock *models.Stock) bool {
	const (
		Period = 20
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

// ComputeVolumeMA is helper method that returns simple moving average of trading volumes for the given period
func ComputeVolumeMA(candles []models.Candle, period int) []float64 {
	if len(candles) < period {
		return utils.EmptySlice[float64]()
	}

	sum := 0.0
	volumeMA := make([]float64, 0, constants.LookBackDays)
	for index := range candles {
		candle := &candles[index]
		sum += float64(candle.Volume)
		if index+1 >= period {
			volumeMA = append(volumeMA, sum/float64(period))
			sum -= float64(candles[index-period+1].Volume)
		}
	}

	return volumeMA
}
