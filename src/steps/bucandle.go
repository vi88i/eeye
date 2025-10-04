package steps

import (
	"eeye/src/models"
	"eeye/src/store"
	"log"
)

func isSolid(candle *models.Candle) bool {
	var (
		openPrice  = candle.Open
		closePrice = candle.Close
		lowPrice   = candle.Low
		highPrice  = candle.High
	)

	if openPrice >= closePrice {
		return false
	}

	var (
		body  = closePrice - openPrice
		upper = highPrice - closePrice
		lower = openPrice - lowPrice
		total = highPrice - lowPrice
	)

	if body >= 0.6*total &&
		upper <= 0.25*body &&
		lower <= 0.25*body {
		log.Printf("Solid Bullish candle: %v\n", candle.Symbol)
		return true
	}

	return false
}

func isHammer(candle *models.Candle) bool {
	var (
		openPrice  = candle.Open
		closePrice = candle.Close
		lowPrice   = candle.Low
		highPrice  = candle.High
	)

	if openPrice >= closePrice {
		return false
	}

	var (
		body  = closePrice - openPrice
		upper = highPrice - closePrice
		lower = openPrice - lowPrice
	)

	if lower < 2*body || upper > 0.25*body || (highPrice-closePrice) > body {
		return false
	}

	log.Printf("Hammer candle: %v\n", candle.Symbol)
	return true
}

func isEngulfing(candle1 *models.Candle, candle2 *models.Candle) bool {
	var (
		openPrice1  = candle1.Open
		closePrice1 = candle1.Close
		openPrice2  = candle2.Open
		closePrice2 = candle2.Close
	)

	if closePrice1 >= openPrice1 || closePrice2 <= openPrice2 {
		return false
	}

	if openPrice2 <= closePrice1 && closePrice2 >= openPrice1 {
		log.Printf("Engulfing pattern: %v\n", candle1.Symbol)
		return true
	}

	return false
}

func isPiercing(candle1 *models.Candle, candle2 *models.Candle) bool {
	var (
		openPrice1  = candle1.Open
		closePrice1 = candle1.Close
		openPrice2  = candle2.Open
		closePrice2 = candle2.Close
	)

	if closePrice1 >= openPrice1 || closePrice2 <= openPrice2 {
		return false
	}

	midpoint := (openPrice1 + closePrice1) / 2
	if openPrice2 < closePrice1 && closePrice2 > midpoint {
		log.Printf("Piercing pattern: %v\n", candle1.Symbol)
		return true
	}

	return false
}

// BullishCandle creates a function that screens for stocks showing bullish
// candlestick patterns. It checks for various bullish patterns including:
// - Solid bullish candles (strong upward momentum)
// - Hammer patterns (potential reversal at support)
// - Engulfing patterns (trend reversal signal)
// - Piercing patterns (bullish reversal after downtrend)
type BullishCandle struct{}

//revive:disable-next-line exported
func (b *BullishCandle) Name() string {
	return "Bullish candle screener"
}

//revive:disable-next-line exported
func (b *BullishCandle) Screen(strategy string, stock *models.Stock) bool {
	const (
		MinPoints = 1
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

	isTwoCandleStickPatternValid := length >= 2
	test := isSolid(&candles[length-1]) ||
		isHammer(&candles[length-1]) ||
		(isTwoCandleStickPatternValid &&
			isEngulfing(&candles[length-2], &candles[length-1])) ||
		(isTwoCandleStickPatternValid &&
			isPiercing(&candles[length-2], &candles[length-1]))

	if !test {
		log.Printf("[%v - %v] test failed: %v\n", strategy, step, stock.Symbol)
	}
	return test
}
