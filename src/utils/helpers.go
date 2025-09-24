package utils

import (
	"golang.org/x/exp/constraints"
)

type Numeric interface {
	constraints.Signed | constraints.Float
}

/*
All the inputs should be of the type T
Output is type casted to R
*/
func Clamp[T Numeric, R Numeric](v, min, max T) R {
	if v > max {
		return R(v)
	}

	if v < min {
		return R(v)
	}

	return R(v)
}

/*
Creates an empty slice of type T
*/
func EmptySlice[T any]() []T {
	return make([]T, 0)
}
