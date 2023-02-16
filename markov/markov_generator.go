package markov

import (
	"diploma/poisson"
	"gonum.org/v1/gonum/mat"
	"math/rand"
)

type Generator struct {
	mu, a     *mat.Dense
	processes []*poisson.ProcessWithIntensityFunc
}

func NewGenerator(mu, a *mat.Dense, processes []*poisson.ProcessWithIntensityFunc) *Generator {
	// todo: sizes
	return &Generator{mu: mu, a: a, processes: processes}
}

func (m *Generator) Generate(n int) []*poisson.Area {
	_, col := m.a.Dims()
	base := rand.Float64()

	areas := make([]*poisson.Area, 0)

	var baseState = 0
	for baseState = 0; baseState < col; baseState++ {
		if m.mu.At(0, baseState) > base {
			break
		}

		base -= m.mu.At(0, baseState)
	}

	areas = append(areas, m.processes[baseState].Generate())

	for i := 1; i < n; i++ {
		//fmt.Println(baseState)
		base = rand.Float64()
		for j := 0; j < col; j++ {
			if m.a.At(baseState, j) > base {
				baseState = j

				break
			}

			base -= m.a.At(baseState, j)
		}

		areas = append(areas, m.processes[baseState].Generate())
	}

	return areas
}
