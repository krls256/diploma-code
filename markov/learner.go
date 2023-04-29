package markov

import (
	"diploma/poisson"
	"gonum.org/v1/gonum/mat"
	"math"
)

type Learner struct {
	StartModel *Model
	Model      *Model
	Frames     []*poisson.Area

	ScalingCoefficients []float64 // sum of alphas (it is 1/c in Nilsson terminology)
	LogProbs            []float64
}

func NewLearner(model *Model, frames []*poisson.Area) *Learner {
	return &Learner{StartModel: model.DeepCopy(), Model: model, Frames: frames}
}

func (l *Learner) CalcB(state int, area *poisson.Area) float64 {
	return PoissonProbability(state, area, l.Model.ObservableProcesses)
}

func (l *Learner) CalcAlphas(useScale bool) (alphas [][]float64, scalingCoefficients []float64) {
	return Forward(useScale, l.Frames, l.Model)
}

func (l *Learner) CalcBetas(useScale bool, scalingCoefficients []float64) [][]float64 {
	return Backward(useScale, scalingCoefficients, l.Frames, l.Model)
}

func (l *Learner) Step() (scalingCoefficients []float64) {
	useScale := true
	_, states := l.Model.BaseDistribution.Dims()
	T := len(l.Frames)

	alphas, scalingCoefficients := l.CalcAlphas(useScale)
	betas := l.CalcBetas(useScale, scalingCoefficients)

	newMu := mat.NewDense(1, states, make([]float64, states))
	newA := mat.NewDense(states, states, make([]float64, states*states))
	newProcesses := make([]poisson.IntensityMap, states)

	for i := 0; i < states; i++ {
		divider := 0.0

		for j := 0; j < states; j++ {
			divider += betas[0][j] * alphas[0][j]
		}

		newMu.Set(0, i, betas[0][i]*alphas[0][i]/divider)
	}

	for i := 0; i < states; i++ {
		for j := 0; j < states; j++ {
			bottom := 0.0
			top := 0.0

			for k := 0; k < states; k++ {
				tmp := 0.0
				for t := 0; t < T-1; t++ {
					tmp += alphas[t][i] * betas[t+1][k] * l.Model.HiddenDistribution.At(i, k) * l.CalcB(k, l.Frames[t+1])
				}

				bottom += tmp
			}

			for t := 0; t < T-1; t++ {
				top += alphas[t][i] * betas[t+1][j] * l.Model.HiddenDistribution.At(i, j) * l.CalcB(j, l.Frames[t+1])
			}

			newA.Set(i, j, top/bottom)
		}
	}

	for i := 0; i < states; i++ {
		newProcesses[i] = poisson.IntensityMap{}
		// in formulas region is j
		for region := range *l.Frames[0] {
			top := 0.0
			bottom := 0.0
			for t := 0; t < T-1; /* WARN not T+1 */ t++ {
				subSum := 0.0
				count := (*l.Frames[t])[region]

				for k := 0; k < states; k++ {
					subSum += l.Model.HiddenDistribution.At(i, k) * betas[t+1][k] * l.CalcB(k, l.Frames[t+1])
				}

				bottom += alphas[t][i] * subSum
				top += alphas[t][i] * subSum * float64(count)

			}

			newProcesses[i][region] = top / bottom
		}
	}

	l.Model.BaseDistribution = newMu
	l.Model.HiddenDistribution = newA
	l.Model.ObservableProcesses = newProcesses

	l.ScalingCoefficients = scalingCoefficients

	return scalingCoefficients
}

func (l *Learner) LogProb() float64 {
	var logProb float64

	for i := 0; i < len(l.ScalingCoefficients); i++ {
		logProb -= math.Log2(1 / l.ScalingCoefficients[i])
	}

	return logProb
}
