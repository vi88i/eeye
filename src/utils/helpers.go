// Package utils provides utility functions and helpers used across the application.
// It includes functions for time formatting, data parsing, type constraints,
// and other common operations.
package utils

import (
	"golang.org/x/exp/constraints"
)

// Numeric is a constraint interface that permits any numeric type that supports
// signed arithmetic or floating-point operations.
type Numeric interface {
	constraints.Signed | constraints.Float
}

// Clamp ensures a value stays within specified bounds by clamping it between minVal and maxVal.
// It accepts any numeric type T for input and returns type R for flexibility in type conversion.
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
// This is a convenience function for creating type-safe empty slices.
func EmptySlice[T any]() []T {
	return make([]T, 0)
}

// Last returns the last value of
func Last[T any](items []T, empty T) T {
	length := len(items)
	if length == 0 {
		return empty
	}

	return items[length-1]
}
