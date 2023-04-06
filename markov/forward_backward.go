package markov

import (
	"diploma/poisson"
)

func Forward(useScale bool, frames []*poisson.Area, model *Model) (alphas [][]float64, scalingCoefficients []float64) {
	_, stateCount := model.Mu.Dims()
	alphas = make([][]float64, len(frames))

	for t := 0; t < len(frames); t++ {
		scaleSum := 0.0
		alphas[t] = make([]float64, stateCount)

		for j := 0; j < stateCount; j++ {
			if t == 0 {
				alphas[0][j] = model.Mu.At(0, j) * PoissonProbability(j, frames[0], model.HiddenProcesses)
			} else {
				alp := 0.0

				for i := 0; i < stateCount; i++ {
					alp += alphas[t-1][i] * model.A.At(i, j) * PoissonProbability(j, frames[t], model.HiddenProcesses)
				}

				alphas[t][j] = alp
			}

			scaleSum += alphas[t][j]
		}

		scalingCoefficients = append(scalingCoefficients, scaleSum)

		if useScale {
			for j := 0; j < stateCount; j++ {
				alphas[t][j] /= scaleSum
			}
		}

	}

	return alphas, scalingCoefficients
}

func Backward(useScale bool, scalingCoefficients []float64, frames []*poisson.Area, model *Model) [][]float64 {
	_, states := model.Mu.Dims()
	betas := make([][]float64, len(frames))

	for t := len(frames) - 1; t >= 0; t-- {
		scaleSum := 0.0
		betas[t] = make([]float64, states)

		for j := 0; j < states; j++ {
			if t == len(frames)-1 {
				betas[t][j] = 1
			} else {
				beta := 0.0
				for i := 0; i < states; i++ {
					beta += model.A.At(j, i) * PoissonProbability(i, frames[t+1], model.HiddenProcesses) * betas[t+1][i]
				}

				betas[t][j] = beta
			}

			scaleSum += betas[t][j]
		}

		if useScale {
			for j := 0; j < states; j++ {
				betas[t][j] /= scalingCoefficients[t]
			}
		}
	}

	return betas
}
