package utils

import (
	"golang.org/x/exp/rand"
	"time"
)

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}

func RandomizeAndNormalizeRows(rows, cols int, mat []float64, randomizeCoefficient float64) []float64 {
	if len(mat) != rows*cols {
		panic("wrong sizes")
	}

	for row := 0; row < rows; row++ {
		sum := 0.0
		for col := 0; col < cols; col++ {
			mat[row*cols+col] += rand.Float64() * randomizeCoefficient
			sum += mat[row*cols+col]
		}

		for col := 0; col < cols; col++ {
			mat[row*cols+col] /= sum
		}
	}

	return mat
}

func RandomStochasticMatrix(rows, cols int) []float64 {
	return RandomizeAndNormalizeRows(rows, cols, make([]float64, rows*cols), 1)
}

func RandomizeMap[K comparable](m map[K]float64, randomizeCoefficient float64) map[K]float64 {
	for key := range m {
		m[key] += rand.Float64() * randomizeCoefficient
	}

	return m
}

func RandomMap[K comparable](keys []K, expectedSum float64, randomizeCoefficient float64) map[K]float64 {
	m := map[K]float64{}

	for _, key := range keys {
		m[key] = rand.Float64() * randomizeCoefficient
	}

	return m
}
