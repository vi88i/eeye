package steps

import (
	"eeye/src/models"
	"eeye/src/store"
	"log"
)

// isSolid checks if a candle is a solid bullish candle.
// A solid bullish candle indicates strong upward momentum with minimal wicks.
// Criteria:
//   - Close must be higher than open (bullish candle)
//   - Body must be at least 60% of total candle range
//   - Upper wick must be <= 25% of body
//   - Lower wick must be <= 25% of body
func isSolid(candle *models.Candle) bool {
	var (
		openPrice  = candle.Open
		closePrice = candle.Close
		lowPrice   = candle.Low
		highPrice  = candle.High
	)

	// Must be bullish (close > open)
	if openPrice >= closePrice {
		return false
	}

	var (
		body  = closePrice - openPrice // Size of candle body
		upper = highPrice - closePrice // Upper wick (above close)
		lower = openPrice - lowPrice   // Lower wick (below open)
		total = highPrice - lowPrice   // Total candle range
	)

	// Check if body is dominant (60%+) with minimal wicks (25% max)
	if body >= 0.6*total &&
		upper <= 0.25*body &&
		lower <= 0.25*body {
		log.Printf("Solid Bullish candle: %v\n", candle.Symbol)
		return true
	}

	return false
}

// isHammer checks if a candle is a hammer pattern.
// A hammer is a bullish reversal pattern typically found at the bottom of a downtrend.
// Criteria:
//   - Close must be higher than open (bullish)
//   - Long lower wick (at least 2x the body size)
//   - Small upper wick (<= 25% of body)
//   - Upper wick must not exceed body size
func isHammer(candle *models.Candle) bool {
	var (
		openPrice  = candle.Open
		closePrice = candle.Close
		lowPrice   = candle.Low
		highPrice  = candle.High
	)

	// Must be bullish (close > open)
	if openPrice >= closePrice {
		return false
	}

	var (
		body  = closePrice - openPrice // Size of candle body
		upper = highPrice - closePrice // Upper wick
		lower = openPrice - lowPrice   // Lower wick (should be long)
	)

	// Validate hammer criteria: long lower wick, small upper wick
	if lower < 2*body || upper > 0.25*body || (highPrice-closePrice) > body {
		return false
	}

	log.Printf("Hammer candle: %v\n", candle.Symbol)
	return true
}

// isEngulfing checks for a bullish engulfing pattern between two candles.
// This is a strong reversal signal where a bullish candle completely engulfs
// the previous bearish candle's body.
// Criteria:
//   - candle1 must be bearish (close < open)
//   - candle2 must be bullish (close > open)
//   - candle2's body must completely engulf candle1's body
func isEngulfing(candle1 *models.Candle, candle2 *models.Candle) bool {
	var (
		openPrice1  = candle1.Open
		closePrice1 = candle1.Close
		openPrice2  = candle2.Open
		closePrice2 = candle2.Close
	)

	// candle1 must be bearish, candle2 must be bullish
	if closePrice1 >= openPrice1 || closePrice2 <= openPrice2 {
		return false
	}

	// candle2 must engulf candle1 completely
	if openPrice2 <= closePrice1 && closePrice2 >= openPrice1 {
		log.Printf("Engulfing pattern: %v\n", candle1.Symbol)
		return true
	}

	return false
}

// isPiercing checks for a piercing pattern between two candles.
// This is a bullish reversal pattern where a bullish candle "pierces" into
// the previous bearish candle's body, closing above its midpoint.
// Criteria:
//   - candle1 must be bearish (close < open)
//   - candle2 must be bullish (close > open)
//   - candle2 opens below candle1's close
//   - candle2 closes above candle1's midpoint
func isPiercing(candle1 *models.Candle, candle2 *models.Candle) bool {
	var (
		openPrice1  = candle1.Open
		closePrice1 = candle1.Close
		openPrice2  = candle2.Open
		closePrice2 = candle2.Close
	)

	// candle1 must be bearish, candle2 must be bullish
	if closePrice1 >= openPrice1 || closePrice2 <= openPrice2 {
		return false
	}

	// Calculate midpoint of candle1's body
	midpoint := (openPrice1 + closePrice1) / 2
	// candle2 opens below candle1's close and closes above midpoint
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
type BullishCandle struct {
	models.StepBaseImpl
}

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

	// Two-candle patterns (engulfing, piercing) require at least 2 candles
	isTwoCandleStickPatternValid := length >= 2
	return b.TruthyCheck(
		strategy,
		step,
		stock,
		func() bool {
			return isSolid(&candles[length-1]) ||
				isHammer(&candles[length-1]) ||
				(isTwoCandleStickPatternValid &&
					isEngulfing(&candles[length-2], &candles[length-1])) ||
				(isTwoCandleStickPatternValid &&
					isPiercing(&candles[length-2], &candles[length-1]))
		},
	)
}
