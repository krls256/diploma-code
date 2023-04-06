package markov

import (
	"diploma/poisson"
	"gitlab.eemi.tech/golang/factgorial/factgorial"
	"math"
)

func PoissonProbability(state int, area *poisson.Area, processes []poisson.IntensityMap) float64 {
	val := 1.0

	for region, count := range *area {
		g := processes[state][region]

		f, err := factgorial.Factorial(count)
		if err != nil {
			panic(err)
		}

		val *= math.Pow(g, float64(count)) / float64(f) * math.Pow(math.E, -g)
	}

	return val
}
