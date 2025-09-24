package steps

import (
	"eeye/src/constants"
	"eeye/src/models"
	"log"
)

func VolumeScreener(
	strategy string,
	stock *models.Stock,
	screen func(currentVolume float64, averageVolume float64) bool,
) func() bool {
	return func() bool {
		const (
			Period = 20
		)

		candles, err := getCachedCandles(stock)
		if err != nil {
			return false
		}

		length := len(candles)
		if length < Period {
			log.Printf("insufficient candles for volume screener: %v\n", stock.Symbol)
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
		test := screen(float64(candles[length-1].Volume), volumeMA[maLength-1])
		if !test {
			log.Printf("[%v] %v failed volume screener test\n", strategy, stock.Symbol)
		}
		return test
	}
}
