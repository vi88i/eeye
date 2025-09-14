package steps

import (
	"eeye/src/constants"
	"eeye/src/models"
	"log"
)

const (
	windowSize = 20
)

// It computes the 20-day MA and applies the screening function to filter stocks.
func VolumeScreener(
	stock *models.Stock,
	screen func(currentVolume float64, averageVolume float64) bool,
) bool {
	candles, err := GetCachedCandles(stock)
	if err != nil {
		return false
	}

	length := len(candles)
	if length < windowSize {
		log.Printf("insufficient candles for volume screener: %v\n", stock.Symbol)
		return false
	}

	sum := 0.0
	volumeMA := make([]float64, 0, constants.LookBackDays)
	for index, candle := range candles {
		sum += float64(candle.Volume)
		if index + 1 >= windowSize {
			volumeMA = append(volumeMA, sum / windowSize)
			sum -= float64(candles[index - windowSize + 1].Volume)
		}
	}

	maLength := len(volumeMA)
	return screen(float64(candles[length - 1].Volume), volumeMA[maLength - 1])
}
