package steps

import (
	"eeye/src/constants"
	"eeye/src/models"
	"eeye/src/store"
	"log"
)

// Volume creates a function that screens stocks based on their trading volume.
// It compares the current volume against the average volume using a custom screening
// function to identify significant volume patterns.
type Volume struct {
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
	if length < Period {
		log.Printf("[%v - %v] insufficient candles: %v\n", strategy, step, stock.Symbol)
		return false
	}

	sum := 0.0
	volumeMA := make([]float64, 0, constants.LookBackDays)
	for index := range candles {
		candle := &candles[index]
		sum += float64(candle.Volume)
		if index+1 >= Period {
			volumeMA = append(volumeMA, sum/Period)
			sum -= float64(candles[index-Period+1].Volume)
		}
	}

	maLength := len(volumeMA)
	test := v.Test(float64(candles[length-1].Volume), volumeMA[maLength-1])
	if !test {
		log.Printf("[%v - %v] test failed: %v\n", strategy, step, stock.Symbol)
	}
	return test
}
