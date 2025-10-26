// Package utils provides utility functions and helpers used across the application.
// It includes functions for time formatting, data parsing, type constraints,
// and other common operations.
package utils

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"

	progressbar "github.com/schollz/progressbar/v3"
	"golang.org/x/exp/constraints"
)

// Numeric is a constraint interface that permits any numeric type that supports
// signed arithmetic or floating-point operations.
type Numeric interface {
	constraints.Signed | constraints.Float
}

// Clamp ensures a value stays within specified bounds by clamping it between minVal and maxVal.
// It accepts any numeric type T for input and returns type R for flexibility in type conversion.
//
// Parameters:
//   - v: The value to clamp
//   - minVal: Minimum allowed value
//   - maxVal: Maximum allowed value
//
// Returns:
//   - The clamped value of type R
//
// Example:
//
//	Clamp[int, int](15, 0, 10) // returns 10
//	Clamp[float64, float64](-5.0, 0.0, 100.0) // returns 0.0
func Clamp[T Numeric, R Numeric](v, minVal, maxVal T) R {
	if v > maxVal {
		return R(maxVal)
	}

	if v < minVal {
		return R(minVal)
	}

	return R(v)
}

// EmptySlice creates and returns an empty slice of the specified type T.
// This is a convenience function for creating type-safe empty slices,
// particularly useful when returning early from functions that return slices.
//
// Returns:
//   - An empty slice of type []T with zero capacity
//
// Example:
//
//	emptyInts := EmptySlice[int]() // returns []int{}
func EmptySlice[T any]() []T {
	return make([]T, 0)
}

// Last returns the last element of a slice, or a default value if the slice is empty.
// This provides a safe way to access the last element without index calculations.
//
// Parameters:
//   - items: The slice to get the last element from
//   - empty: Default value to return if slice is empty
//
// Returns:
//   - The last element if slice is non-empty, otherwise the default value
//
// Example:
//
//	Last([]int{1, 2, 3}, 0) // returns 3
//	Last([]int{}, 0) // returns 0
func Last[T any](items []T, empty T) T {
	length := len(items)
	if length == 0 {
		return empty
	}

	return items[length-1]
}

// SelectFields creates an array of maps containing only specified fields from struct items.
// Uses reflection to extract field values from structs dynamically.
//
// Parameters:
//   - items: Slice of struct items to extract fields from
//   - fields: Names of fields to include in the output maps
//
// Returns:
//   - Slice of maps where each map contains only the specified fields
//
// Panics:
//   - If T is not a struct type
//   - If any field name in fields doesn't exist in the struct
//
// Example:
//
//	type Person struct { Name string; Age int; Email string }
//	people := []Person{{"Alice", 30, "alice@example.com"}}
//	SelectFields(people, []string{"Name", "Age"})
//	// returns []map[string]any{{"Name": "Alice", "Age": 30}}
func SelectFields[T any](items []T, fields []string) []map[string]any {
	// Get type information from a zero value of T
	// (cannot use reflect.TypeOf(T) directly as T is a type parameter)
	var zero T
	t := reflect.TypeOf(zero)
	if t.Kind() != reflect.Struct {
		panic("SelectFields: T must of struct type")
	}

	// Build list of available fields in the struct
	tFields := make([]string, 0, t.NumField())
	for i := range t.NumField() {
		tFields = append(tFields, t.Field(i).Name)
	}

	// Validate that all requested fields exist
	for i := range fields {
		key := fields[i]
		if !slices.Contains(tFields, key) {
			panic(fmt.Sprintf("SelectFields: %v field not found", key))
		}
	}

	// Extract requested fields from each item
	ret := make([]map[string]any, 0, len(items))
	for i := range items {
		v := reflect.ValueOf(items[i])
		m := map[string]any{}
		for j := range fields {
			key := fields[j]
			m[key] = v.FieldByName(key).Interface()
		}
		ret = append(ret, m)
	}

	return ret
}

// Map transforms each element of a slice using the provided function.
// This is a functional programming utility similar to map in other languages.
//
// Parameters:
//   - items: Source slice to transform
//   - fn: Transformation function that takes an element of type T and returns type U
//
// Returns:
//   - New slice containing transformed elements of type U
//
// Example:
//
//	numbers := []int{1, 2, 3}
//	doubled := Map(numbers, func(n int) int { return n * 2 })
//	// doubled = []int{2, 4, 6}
func Map[T any, U any](items []T, fn func(T) U) []U {
	res := make([]U, 0, len(items))

	for i := range items {
		res = append(res, fn(items[i]))
	}

	return res
}

// PadLeft pads a slice at the beginning with the given filler value until it reaches the desired length.
// If the slice is already longer than or equal to total, returns the original slice unchanged.
//
// Parameters:
//   - items: The slice to pad
//   - total: Desired total length
//   - fill: Value to use for padding
//
// Returns:
//   - New slice padded to the specified length, or original slice if already long enough
//
// Example:
//
//	PadLeft([]int{3, 4}, 5, 0) // returns []int{0, 0, 0, 3, 4}
//	PadLeft([]int{1, 2, 3}, 2, 0) // returns []int{1, 2, 3} (unchanged)
func PadLeft[T any](items []T, total int, fill T) []T {
	currentLength := len(items)
	if currentLength >= total {
		return items
	}

	// Calculate how many padding elements needed
	rem := total - currentLength
	res := make([]T, 0, rem)
	for range rem {
		res = append(res, fill)
	}

	return append(res, items...)
}

// Round2 rounds a float64 value to two decimal places.
// Uses string formatting for precise decimal rounding.
//
// Parameters:
//   - num: The number to round
//
// Returns:
//   - The number rounded to 2 decimal places
//
// Example:
//
//	Round2(3.14159) // returns 3.14
//	Round2(2.996) // returns 3.00
func Round2(num float64) float64 {
	rounded, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", num), 64)
	return rounded
}

// Reduce applies a reducer function to accumulate a single value from a slice.
// This is a functional programming utility similar to reduce/fold in other languages.
//
// Parameters:
//   - items: The slice to reduce
//   - reducer: Function that takes (accumulator, currentValue, index) and returns new accumulator
//   - initial: Initial value for the accumulator
//
// Returns:
//   - The final accumulated value
//
// Example:
//
//	sum := Reduce([]int{1, 2, 3, 4}, func(acc, val, _ int) int {
//	    return acc + val
//	}, 0) // returns 10
func Reduce[T any](items []T, reducer func(T, T, int) T, initial T) T {
	acc := initial
	for i := range items {
		acc = reducer(acc, items[i], i)
	}
	return acc
}

// Filter creates a new slice containing only elements that satisfy the condition.
// This is a functional programming utility similar to filter in other languages.
//
// Parameters:
//   - items: The slice to filter
//   - condition: Function that takes (element, index) and returns true to keep the element
//
// Returns:
//   - New slice containing only elements where condition returned true
//
// Example:
//
//	evens := Filter([]int{1, 2, 3, 4, 5}, func(n, _ int) bool {
//	    return n % 2 == 0
//	}) // returns []int{2, 4}
func Filter[T any](items []T, condition func(T, int) bool) []T {
	res := []T{}
	for i := range items {
		if condition(items[i], i) {
			res = append(res, items[i])
		}
	}
	return res
}

// GetProgressTracker creates and returns a configured progress bar for tracking operation progress.
// The progress bar is displayed in the terminal with a custom green theme and color-coded output.
//
// Parameters:
//   - num: Total number of items/steps to track (sets the progress bar's maximum value)
//   - description: Text description displayed alongside the progress bar (e.g., "Analyzing stocks...")
//
// Returns:
//   - *progressbar.ProgressBar: A configured progress bar instance ready for tracking
//
// Example usage:
//
//	bar := GetProgressTracker(100, "Processing files...")
//	for i := 0; i < 100; i++ {
//	    // Do work...
//	    bar.Add(1)
//	}
//
// The progress bar displays with format: "[description] [===>    ] 50/100"
func GetProgressTracker(num int, description string) *progressbar.ProgressBar {
	return progressbar.NewOptions(num,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(25),
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
}
