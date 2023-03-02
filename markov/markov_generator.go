package markov

import (
	"diploma/poisson"
	"gonum.org/v1/gonum/mat"
	"math/rand"
)

type Model struct {
	mu, a     *mat.Dense
	processes []*poisson.ProcessWithIntensityFunc
}

func NewModel(mu, a *mat.Dense, processes []*poisson.ProcessWithIntensityFunc) *Model {
	muRows, muCols := mu.Dims()
	if muRows != 1 {
		panic("wrong sizes")
	}

	aRows, aCols := a.Dims()
	if muCols != aCols || aCols != aRows || aRows != len(processes) {
		panic("wrong sizes")
	}
	return &Model{mu: mu, a: a, processes: processes}
}

func (m *Model) Generate(n int) (areas []*poisson.Area, stateChain []int) {
	_, col := m.a.Dims()
	base := rand.Float64()

	var baseState = 0
	for baseState = 0; baseState < col; baseState++ {
		if m.mu.At(0, baseState) > base {
			break
		}

		base -= m.mu.At(0, baseState)
	}

	areas = append(areas, m.processes[baseState].Generate())
	stateChain = append(stateChain, baseState)

	for i := 1; i < n; i++ {
		base = rand.Float64()

		for j := 0; j < col; j++ {
			if m.a.At(baseState, j) > base {
				baseState = j

				break
			}

			base -= m.a.At(baseState, j)
		}

		areas = append(areas, m.processes[baseState].Generate())
		stateChain = append(stateChain, baseState)
	}

	return
}
