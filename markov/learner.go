package markov

import (
	"diploma/poisson"
	"gonum.org/v1/gonum/mat"
	"math"
)

type Learner struct {
	StartModel           *Model
	Model                *Model
	ObservationSequences [][]*poisson.Area

	History []*Model

	ScalingCoefficients []float64 // sum of alphas (it is 1/c in Nilsson terminology)
	LogProbs            []float64
}

func NewLearner(model *Model, observationSequences [][]*poisson.Area) *Learner {
	return &Learner{StartModel: model.DeepCopy(), Model: model, ObservationSequences: observationSequences}
}

func (l *Learner) CalcB(state int, area *poisson.Area) float64 {
	return PoissonProbability(state, area, l.Model.ObservableProcesses)
}

func (l *Learner) Step() (scalingCoefficients []float64) {
	useScale := true
	_, states := l.Model.BaseDistribution.Dims()

	obsSeq := len(l.ObservationSequences)

	// [[1, 0], [0, 1], [1, 0], ...] <- mu for each observation
	muArr := []*mat.Dense{}
	aTopArr := []*mat.Dense{}
	aBottomArr := []*mat.Dense{}
	processTopArr := [][]poisson.IntensityMap{}
	processBottomArr := [][]poisson.IntensityMap{}

	for o := 0; o < obsSeq; o++ {
		frames := l.ObservationSequences[o]
		T := len(frames)

		alphas, scalingCoefficients := Forward(useScale, frames, l.Model)
		betas := Backward(useScale, scalingCoefficients, frames, l.Model)

		newMu := mat.NewDense(1, states, make([]float64, states))
		newATop := mat.NewDense(states, states, make([]float64, states*states))
		newABottom := mat.NewDense(states, states, make([]float64, states*states))
		newProcessesTop := make([]poisson.IntensityMap, states)
		newProcessesBottom := make([]poisson.IntensityMap, states)

		for i := 0; i < states; i++ {
			divider := 0.0

			for j := 0; j < states; j++ {
				divider += betas[0][j] * alphas[0][j]
			}

			newMu.Set(0, i, betas[0][i]*alphas[0][i]/divider)
		}

		muArr = append(muArr, newMu)

		for i := 0; i < states; i++ {
			for j := 0; j < states; j++ {
				bottom := 0.0
				top := 0.0

				for k := 0; k < states; k++ {
					tmp := 0.0
					for t := 0; t < T-1; t++ {
						tmp += alphas[t][i] * betas[t+1][k] * l.Model.HiddenDistribution.At(i, k) * l.CalcB(k, frames[t+1])
					}

					bottom += tmp
				}

				for t := 0; t < T-1; t++ {
					top += alphas[t][i] * betas[t+1][j] * l.Model.HiddenDistribution.At(i, j) * l.CalcB(j, frames[t+1])
				}

				newATop.Set(i, j, top)
				newABottom.Set(i, j, bottom)
			}
		}

		aTopArr = append(aTopArr, newATop)
		aBottomArr = append(aBottomArr, newABottom)

		for i := 0; i < states; i++ {
			newProcessesTop[i] = poisson.IntensityMap{}
			newProcessesBottom[i] = poisson.IntensityMap{}

			// in formulas region is j
			for region := range *frames[0] {
				top := 0.0
				bottom := 0.0
				for t := 0; t < T-1; t++ {
					subSum := 0.0
					count := (*frames[t])[region]

					for k := 0; k < states; k++ {
						subSum += l.Model.HiddenDistribution.At(i, k) * betas[t+1][k] * l.CalcB(k, frames[t+1])
					}

					bottom += alphas[t][i] * subSum
					top += alphas[t][i] * subSum * float64(count)
				}

				top += alphas[T-1][i] * float64((*frames[T-1])[region])
				bottom += alphas[T-1][i]

				newProcessesTop[i][region] = top
				newProcessesBottom[i][region] = bottom
			}
		}

		processTopArr = append(processTopArr, newProcessesTop)
		processBottomArr = append(processBottomArr, newProcessesBottom)
	}

	newMu := mat.NewDense(1, states, make([]float64, states))

	newA := mat.NewDense(states, states, make([]float64, states*states))
	newATop := mat.NewDense(states, states, make([]float64, states*states))
	newABottom := mat.NewDense(states, states, make([]float64, states*states))

	newProcesses := make([]poisson.IntensityMap, states)
	newProcessesTop := make([]poisson.IntensityMap, states)
	newProcessesBottom := make([]poisson.IntensityMap, states)

	for i := 0; i < obsSeq; i++ {
		newMu.Add(newMu, muArr[i])

		newATop.Add(newATop, aTopArr[i])
		newABottom.Add(newABottom, aBottomArr[i])

		for j := 0; j < len(newProcesses); j++ {
			if newProcessesTop[j] == nil {
				newProcessesTop[j] = poisson.IntensityMap{}
			}

			if newProcessesBottom[j] == nil {
				newProcessesBottom[j] = poisson.IntensityMap{}
			}

			newProcessesTop[j].Add(newProcessesTop[j], processTopArr[i][j])
			newProcessesBottom[j].Add(newProcessesBottom[j], processBottomArr[i][j])
		}
	}

	for i := 0; i < states; i++ {
		newMu.Set(0, i, newMu.At(0, i)/float64(obsSeq))

		for j := 0; j < states; j++ {
			newA.Set(i, j, newATop.At(i, j)/newABottom.At(i, j))
		}

		for reg := range newProcessesTop[i] {
			if newProcesses[i] == nil {
				newProcesses[i] = poisson.IntensityMap{}
			}

			newProcesses[i][reg] = newProcessesTop[i][reg] / newProcessesBottom[i][reg]
		}
	}

	l.History = append(l.History, NewModel(l.Model.BaseDistribution, l.Model.HiddenDistribution, l.Model.ObservableProcesses))

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
