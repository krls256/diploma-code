package markov

import (
	"bytes"
	"diploma/poisson"
	"diploma/utils"
	"fmt"
	"github.com/samber/lo"
	"gonum.org/v1/gonum/mat"
	"math/rand"
	"strings"
)

type Model struct {
	Mu, A     *mat.Dense
	Processes []poisson.IntensityMap
}

func (m *Model) String() string {
	buf := bytes.Buffer{}

	muStr := lo.Map(m.Mu.RawMatrix().Data, func(item float64, index int) string {
		return utils.PrecisionString(item)
	})
	buf.WriteString(fmt.Sprintf("Mu: (%v)\n", strings.Join(muStr, ", ")))

	rows, _ := m.A.Dims()

	buf.WriteString("A: \n")
	for i := 0; i < rows; i++ {
		row := m.A.RawRowView(i)

		rowStr := lo.Map(row, func(item float64, index int) string {
			return utils.PrecisionString(item)
		})

		buf.WriteString(fmt.Sprintf("|%v| \n", strings.Join(rowStr, ", ")))
	}

	for i := 0; i < len(m.Processes); i++ {
		buf.WriteString(fmt.Sprintf("process %v:\n", i+1))
		buf.WriteString(m.Processes[i].String())
		if i != len(m.Processes)-1 {
			buf.WriteString("\n")
		}
	}

	return buf.String()
}

func NewModel(mu, a *mat.Dense, processes []poisson.IntensityMap) *Model {
	muRows, muCols := mu.Dims()
	if muRows != 1 {
		panic("wrong sizes")
	}

	aRows, aCols := a.Dims()
	if muCols != aCols || aCols != aRows || aRows != len(processes) {
		panic("wrong sizes")
	}
	return &Model{Mu: mu, A: a, Processes: processes}
}

func (m *Model) Generate(n int) (areas []*poisson.Area, stateChain []int) {
	_, col := m.A.Dims()
	base := rand.Float64()

	var baseState = 0
	for baseState = 0; baseState < col; baseState++ {
		if m.Mu.At(0, baseState) > base {
			break
		}

		base -= m.Mu.At(0, baseState)
	}

	areas = append(areas, m.Processes[baseState].Generate())
	stateChain = append(stateChain, baseState)

	for i := 1; i < n; i++ {
		base = rand.Float64()

		for j := 0; j < col; j++ {
			if m.A.At(baseState, j) > base {
				baseState = j

				break
			}

			base -= m.A.At(baseState, j)
		}

		areas = append(areas, m.Processes[baseState].Generate())
		stateChain = append(stateChain, baseState)
	}

	return
}
