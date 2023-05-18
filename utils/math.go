package utils

import (
	"github.com/samber/lo"
	"golang.org/x/exp/constraints"
	"math"
)

func Max[T constraints.Ordered](sl []T) T {
	if len(sl) == 0 {
		return *new(T)
	}

	max := sl[0]

	for i := 1; i < len(sl); i++ {
		if max < sl[i] {
			max = sl[i]
		}
	}

	return max
}

func MaxIndex[T constraints.Ordered](sl []T) int {
	if len(sl) == 0 {
		return -1
	}

	maxIndex := 0
	max := sl[0]

	for i := 1; i < len(sl); i++ {
		if max < sl[i] {
			max = sl[i]
			maxIndex = i
		}
	}

	return maxIndex
}

func Scale[T constraints.Float](sl []T) []T {
	sum := lo.Sum(sl)

	return lo.Map(sl, func(item T, index int) T {
		return item / sum
	})
}

func CountDiff[T comparable](a, b []T) (diff int) {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			diff++
		}
	}

	ldiff := int(math.Abs(float64(len(a) - len(b))))

	diff += ldiff

	return diff
}
