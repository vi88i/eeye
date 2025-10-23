package steps

import (
	"eeye/src/models"
	"eeye/src/store"
	"eeye/src/utils"
	"log"
	"math"
	"slices"
)

// LiquidityLevels screens stocks based on Support & Resistance (S&R) analysis.
// Identifies price levels where the stock has historically found support or faced resistance,
// which often act as significant price zones for future trading decisions.
type LiquidityLevels struct {
	models.StepBaseImpl
	// Window is the lookback period on each side to identify local peaks/troughs
	Window int
	// Tolerance is the percentage difference to cluster nearby levels (e.g., 0.02 = 2%)
	Tolerance float64
	// Strength is the minimum number of touches required for a level to be significant
	Strength int
	// Test receives candles and identified support/resistance levels for screening.
	// Parameters:
	//   - candles: Historical price data
	//   - supports: Identified support levels (price floors)
	//   - resistances: Identified resistance levels (price ceilings)
	// Returns true if the stock passes the screening test.
	Test func(candles []models.Candle, supports []float64, resistances []float64) bool
}

//revive:disable-next-line exported
func (s *LiquidityLevels) Name() string {
	return "Liquidity levels"
}

//revive:disable-next-line exported
func (s *LiquidityLevels) Screen(strategy string, stock *models.Stock) bool {
	step := s.Name()

	candles, err := store.Get(stock)
	if err != nil {
		return false
	}

	if s.Window <= 0 {
		log.Printf("[LiquidityLevels] window size %v is not valid, should be > 0\n", s.Window)
		return false
	}

	if s.Strength <= 0 {
		log.Printf("[LiquidityLevels] strength %v is not valid, should be > 0\n", s.Strength)
		return false
	}

	if s.Tolerance <= 0 {
		log.Printf("[LiquidityLevels] tolerance %v is not valid, should be > 0\n", s.Tolerance)
		return false
	}

	supports, resistances := GetLiquidityLevels(candles, s.Window, s.Tolerance, s.Strength)

	return s.TruthyCheck(
		strategy,
		step,
		stock,
		func() bool {
			return s.Test(candles, supports, resistances)
		},
	)
}

// getLocalPeaksAndTroughs returns the local maxima and minima
func getLocalPeaksAndTroughs(candles []models.Candle, window int) ([]float64, []float64) {
	var (
		peaks        = []float64{}
		troughs      = []float64{}
		numOfCandles = len(candles)
	)

	for i := range candles {
		left := i - window
		right := i + window + 1
		if left >= 0 && right <= numOfCandles {
			currentHigh := candles[i].High
			highs := utils.Map(
				candles[left:right],
				func(candle models.Candle) float64 {
					return candle.High
				},
			)

			currentLow := candles[i].Low
			lows := utils.Map(
				candles[left:right],
				func(candle models.Candle) float64 {
					return candle.Low
				},
			)

			if slices.Max(highs) == currentHigh {
				peaks = append(peaks, currentHigh)
			}

			if slices.Min(lows) == currentLow {
				troughs = append(troughs, currentLow)
			}
		}
	}

	slices.Sort(peaks)
	slices.Sort(troughs)
	return peaks, troughs
}

// getLevels uses a simple mean-based clustering algorithm to identify significant price levels.
// Prices are grouped into clusters based on tolerance (percentage difference from cluster mean).
// Only clusters with at least 'strength' number of prices are considered significant.
//
// Parameters:
//   - prices: Sorted list of peak or trough prices
//   - tolerance: Maximum percentage difference from cluster mean to add a price
//   - strength: Minimum cluster size to be considered a significant level
//
// Returns:
//   - Slice of mean prices for each significant cluster
func getLevels(prices []float64, tolerance float64, strength int) []float64 {
	clusters := [][]float64{}
	for i := range prices {
		// Initialize first cluster
		if len(clusters) == 0 {
			clusters = append(clusters, []float64{prices[i]})
			continue
		}

		// Calculate mean of the last cluster
		lastClusterIndex := len(clusters) - 1
		mean := utils.Reduce(
			clusters[lastClusterIndex],
			func(acc float64, currentValue float64, _ int) float64 {
				return acc + currentValue
			},
			0.0,
		) / float64(len(clusters[lastClusterIndex]))

		// Check if current price belongs to last cluster or starts a new one
		if mean != 0 {
			diff := math.Abs(prices[i]-mean) / mean
			if diff > tolerance {
				// Price is too far from cluster mean, start new cluster
				clusters = append(clusters, []float64{prices[i]})
			} else {
				// Price is close enough, add to existing cluster
				clusters[lastClusterIndex] = append(clusters[lastClusterIndex], prices[i])
			}
		} else {
			// Edge case: mean is zero, create new cluster
			clusters = append(clusters, []float64{prices[i]})
		}
	}

	// Filter clusters by strength and calculate mean price for each
	return utils.Map(
		utils.Filter(
			clusters,
			func(cluster []float64, _ int) bool {
				// Only keep clusters with enough touches
				return len(cluster) >= strength
			},
		),
		func(cluster []float64) float64 {
			// Calculate mean price of the cluster
			totalItems := len(cluster)
			sum := utils.Reduce(
				cluster,
				func(acc float64, v float64, _ int) float64 {
					return acc + v
				},
				0.0,
			)
			return sum / float64(totalItems)
		},
	)
}

// GetLiquidityLevels identifies significant support and resistance levels.
// Algorithm:
//  1. Find local peaks and troughs using a rolling window
//  2. Cluster nearby levels using tolerance-based algorithm
//  3. Filter clusters by strength (minimum number of touches)
//  4. Return mean price of each significant cluster
//
// Parameters:
//   - candles: Historical price data
//   - window: Number of candles on each side for peak/trough detection
//   - tolerance: Max percentage difference to group levels (e.g., 0.02 = 2%)
//   - strength: Minimum touches for a level to be considered significant
//
// Returns:
//   - supports: Slice of support levels (price floors)
//   - resistances: Slice of resistance levels (price ceilings)
func GetLiquidityLevels(
	candles []models.Candle,
	window int,
	tolerance float64,
	strength int,
) ([]float64, []float64) {
	peaks, troughs := getLocalPeaksAndTroughs(candles, window)
	supports := getLevels(troughs, tolerance, strength)
	resistance := getLevels(peaks, tolerance, strength)
	return supports, resistance
}
