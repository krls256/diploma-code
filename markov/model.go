package markov

import (
	"bytes"
	"diploma/poisson"
	"diploma/utils"
	"fmt"
	"github.com/samber/lo"
	"gonum.org/v1/gonum/mat"
	"image"
	"image/png"
	"math/rand"
	"os"
	"strings"
)

type Model struct {
	BaseDistribution    *mat.Dense
	HiddenDistribution  *mat.Dense
	ObservableProcesses []poisson.IntensityMap
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
	return &Model{BaseDistribution: mu, HiddenDistribution: a, ObservableProcesses: processes}
}

func (m *Model) Generate(n int) *ResultChain {
	rc := &ResultChain{}
	_, col := m.HiddenDistribution.Dims()
	base := rand.Float64()

	var baseState = 0
	for baseState = 0; baseState < col; baseState++ {
		if m.BaseDistribution.At(0, baseState) > base {
			break
		}

		base -= m.BaseDistribution.At(0, baseState)
	}

	rc.Frames = append(rc.Frames, m.ObservableProcesses[baseState].Generate())
	rc.StateChain = append(rc.StateChain, baseState)

	for i := 1; i < n; i++ {
		base = rand.Float64()

		for j := 0; j < col; j++ {
			if m.HiddenDistribution.At(baseState, j) > base {
				baseState = j

				break
			}

			base -= m.HiddenDistribution.At(baseState, j)
		}

		rc.Frames = append(rc.Frames, m.ObservableProcesses[baseState].Generate())
		rc.StateChain = append(rc.StateChain, baseState)
	}

	return rc
}

func (m *Model) String() string {
	buf := bytes.Buffer{}

	muStr := lo.Map(m.BaseDistribution.RawMatrix().Data, func(item float64, index int) string {
		return utils.PrecisionString(item)
	})
	buf.WriteString(fmt.Sprintf("BaseDistribution: (%v)\n", strings.Join(muStr, ", ")))

	rows, _ := m.HiddenDistribution.Dims()

	buf.WriteString("HiddenDistribution: \n")
	for i := 0; i < rows; i++ {
		row := m.HiddenDistribution.RawRowView(i)

		rowStr := lo.Map(row, func(item float64, index int) string {
			return utils.PrecisionString(item)
		})

		buf.WriteString(fmt.Sprintf("|%v| \n", strings.Join(rowStr, ", ")))
	}

	for i := 0; i < len(m.ObservableProcesses); i++ {
		buf.WriteString(fmt.Sprintf("process %v:\n", i+1))
		buf.WriteString(m.ObservableProcesses[i].String())
		if i != len(m.ObservableProcesses)-1 {
			buf.WriteString("\n")
		}
	}

	return buf.String()
}

func (m *Model) Draw(filename string) {
	imgs := []image.Image{}

	for i, p := range m.ObservableProcesses {
		imgs = append(imgs, poisson.DrawIntensityMap(p, i+1))
	}

	img := utils.HorizontalJoinImage(imgs[0], imgs[1:]...)

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}

func (m *Model) DeepCopy() *Model {
	nm := &Model{
		BaseDistribution:   &mat.Dense{},
		HiddenDistribution: &mat.Dense{},
	}

	nm.BaseDistribution.Copy(m.BaseDistribution)
	nm.HiddenDistribution.Copy(m.HiddenDistribution)

	for i := 0; i < len(m.ObservableProcesses); i++ {
		nm.ObservableProcesses = append(nm.ObservableProcesses, utils.SimpleMapDeepCopy(m.ObservableProcesses[i]))
	}

	return nm
}
