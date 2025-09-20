package steps

import (
	"eeye/src/models"
	"log"
)

func isSolid(candle *models.Candle) bool {
	var (
		open  = candle.Open
		close = candle.Close
		low   = candle.Low
		high  = candle.High
	)

	if open >= close {
		return false
	}

	var (
		body  = close - open
		upper = high - close
		lower = open - low
		total = high - low
	)

	if body >= 0.7*total &&
		upper <= 0.15*body &&
		lower <= 0.15*body {
		log.Printf("Solid Bullish candle: %v\n", candle.Symbol)
		return true
	}

	return false
}

func isHammer(candle *models.Candle) bool {
	var (
		open  = candle.Open
		close = candle.Close
		low   = candle.Low
		high  = candle.High
	)

	if open >= close {
		return false
	}

	var (
		body  = close - open
		upper = high - close
		lower = open - low
	)

	if lower < 2*body || upper > 0.25*body || (high-close) > body {
		return false
	}

	log.Printf("Hammer candle: %v\n", candle.Symbol)
	return true
}

func isEngulfing(candle1 *models.Candle, candle2 *models.Candle) bool {
	var (
		open1  = candle1.Open
		close1 = candle1.Close
		open2  = candle2.Open
		close2 = candle2.Close
	)

	if close1 >= open1 || close2 <= open2 {
		return false
	}

	if open2 <= close1 && close2 >= open1 {
		log.Printf("Engulfing pattern: %v\n", candle1.Symbol)
		return true
	}

	return false
}

func isPiercing(candle1 *models.Candle, candle2 *models.Candle) bool {
	var (
		open1  = candle1.Open
		close1 = candle1.Close
		open2  = candle2.Open
		close2 = candle2.Close
	)

	if close1 >= open1 || close2 <= open2 {
		return false
	}

	midpoint := (open1 + close1) / 2
	if open2 < close1 && close2 > midpoint {
		log.Printf("Piercing pattern: %v\n", candle1.Symbol)
		return true
	}

	return false
}

func BullishCandleScreener(
	strategy string,
	stock *models.Stock,
) func() bool {
	return func() bool {
		candles, err := GetCachedCandles(stock)
		if err != nil {
			return false
		}

		length := len(candles)
		if length <= 0 {
			log.Printf("insufficient candles for candle pattern detection: %v\n", stock.Symbol)
			return false
		}

		isTwoCandleStickPatternValid := length >= 2
		test := isSolid(&candles[length-1]) ||
			isHammer(&candles[length-1]) ||
			(isTwoCandleStickPatternValid &&
				isEngulfing(&candles[length-2], &candles[length-1])) ||
			(isTwoCandleStickPatternValid &&
				isPiercing(&candles[length-2], &candles[length-1]))

		if !test {
			log.Printf("[%v] %v failed bullish candle test\n", strategy, stock.Symbol)
		}

		return test
	}
}
