package utils

import (
	"golang.org/x/exp/constraints"
)

type Numeric interface {
	constraints.Signed | constraints.Float
}

func Clamp[T Numeric, R Numeric](v, min, max T) R {
	if v > max {
		return R(v)
	}

	if v < min {
		return R(v)
	}

	return R(v)
}

func EmptySlice[T any]() []T {
	return make([]T, 0)
}
