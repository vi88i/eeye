// Package utils provides utility functions and helpers used across the application.
// It includes functions for time formatting, data parsing, type constraints,
// and other common operations.
package utils

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"

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

// SelectFields creates an array of map[string]any from the items array
// with only fields specified in the fields array
func SelectFields[T any](items []T, fields []string) []map[string]any {
	// We cannot do reflect.TypeOf(T) because T is a type not a variable (it expects a value)
	// So trick is to create a zero value of T
	var zero T
	t := reflect.TypeOf(zero)
	if t.Kind() != reflect.Struct {
		panic("SelectFields: T must of struct type")
	}

	tFields := make([]string, 0, t.NumField())
	for i := range t.NumField() {
		tFields = append(tFields, t.Field(i).Name)
	}

	for i := range fields {
		key := fields[i]
		if !slices.Contains(tFields, key) {
			panic(fmt.Sprintf("SelectFields: %v field not found", key))
		}
	}

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

// Map helps to map an array of items to target type
func Map[T any, U any](items []T, fn func(T) U) []U {
	res := make([]U, 0, len(items))

	for i := range items {
		res = append(res, fn(items[i]))
	}

	return res
}

// PadLeft pads a slice at the beginning with the given filler value
func PadLeft[T any](items []T, total int, fill T) []T {
	currentLength := len(items)
	if currentLength >= total {
		return items
	}

	rem := total - currentLength
	res := make([]T, 0, rem)
	for range rem {
		res = append(res, fill)
	}

	return append(res, items...)
}

// Round2 rounds off the float64 to two decimal places
func Round2(num float64) float64 {
	rounded, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", num), 64)
	return rounded
}

// Reduce helps to reduce the slice using the reducer function
func Reduce[T any](items []T, reducer func(T, T, int) T, initial T) T {
	acc := initial
	for i := range items {
		acc = reducer(acc, items[i], i)
	}
	return acc
}

// Filter helps to filter the slice using the condition function
func Filter[T any](items []T, condition func(T, int) bool) []T {
	res := []T{}
	for i := range items {
		if condition(items[i], i) {
			res = append(res, items[i])
		}
	}
	return res
}
