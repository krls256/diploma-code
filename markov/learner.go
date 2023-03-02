package markov

import (
	"diploma/poisson"
	"gitlab.eemi.tech/golang/factgorial/factgorial"
	"gonum.org/v1/gonum/mat"
	"math"
)

type Learner struct {
	Mu, A     *mat.Dense
	Processes []map[poisson.Region]float64
	frames    []*poisson.Area
}

func NewLearner(mu, a *mat.Dense, processes []map[poisson.Region]float64, frames []*poisson.Area) *Learner {
	return &Learner{A: a, Mu: mu, Processes: processes, frames: frames}
}

func (l *Learner) CalcB(state int, area *poisson.Area) float64 {
	val := 1.0

	for region, count := range *area {
		g := l.Processes[state][region]

		f, err := factgorial.Factorial(count)
		if err != nil {
			panic(err)
		}

		val *= math.Pow(g, float64(count)) / float64(f) * math.Pow(math.E, -g)
	}

	//if val == 0 {
	//	fmt.Println(l.Processes)
	//	//for region, count := range *area {
	//	//	g := l.Processes[state][region]
	//	//
	//	//	fmt.Println("dd", g, count, math.Pow(g, float64(count)))
	//	//}
	//}

	return val
}

func (l *Learner) CalcAlphas() [][]float64 {
	_, states := l.Mu.Dims()
	alphas := make([][]float64, len(l.frames))

	//alphas[0] = make([]float64, states)
	//for i := 0; i < states; i++ {
	//	alphas[0][i] = l.Mu.At(0, i) * l.CalcB(i, l.frames[0])
	//}

	for t := 0; t < len(l.frames); t++ {
		scaleSum := 0.0
		alphas[t] = make([]float64, states)

		for j := 0; j < states; j++ {
			if t == 0 {
				alphas[0][j] = l.Mu.At(0, j) * l.CalcB(j, l.frames[0])
			} else {
				alp := 0.0

				for i := 0; i < states; i++ {
					alp += alphas[t-1][i] * l.A.At(i, j) * l.CalcB(j, l.frames[t])

					//fmt.Println("ddd", alphas[t-1][i], l.A.At(i, j), l.CalcB(j, l.frames[t]))

				}

				alphas[t][j] = alp
			}

			scaleSum += alphas[t][j]
		}

		//fmt.Println("scaleSum", scaleSum)
		//fmt.Println("alphas[t]", alphas[t])
		for j := 0; j < states; j++ {
			alphas[t][j] /= scaleSum
		}

		//fmt.Println("alphas[t] new", alphas[t])
	}

	return alphas
}

func (l *Learner) CalcBetas() [][]float64 {
	_, states := l.Mu.Dims()
	betas := make([][]float64, len(l.frames))

	//betas[len(l.frames)-1] = make([]float64, states)
	//for i := 0; i < states; i++ {
	//	betas[len(l.frames)-1][i] = 1
	//}

	for t := len(l.frames) - 1; t >= 0; t-- {
		scaleSum := 0.0
		betas[t] = make([]float64, states)

		for j := 0; j < states; j++ {
			if t == len(l.frames)-1 {
				betas[t][j] = 1
			} else {
				beta := 0.0
				for i := 0; i < states; i++ {
					beta += l.A.At(j, i) * l.CalcB(i, l.frames[t+1]) * betas[t+1][i]
				}

				betas[t][j] = beta
			}

			scaleSum += betas[t][j]
		}

		for j := 0; j < states; j++ {
			betas[t][j] /= scaleSum
		}
	}

	return betas
}

func (l *Learner) Step() {
	_, states := l.Mu.Dims()
	T := len(l.frames)
	alphas := l.CalcAlphas()
	//fmt.Println("alphas", alphas[:5])
	betas := l.CalcBetas()
	//fmt.Println("betas", betas[:5])

	newMu := mat.NewDense(1, states, make([]float64, states))
	newA := mat.NewDense(states, states, make([]float64, states*states))
	newProcesses := make([]map[poisson.Region]float64, states)

	for i := 0; i < states; i++ {
		divider := 0.0

		for j := 0; j < states; j++ {
			divider += betas[0][j] * alphas[0][j]
		}

		//fmt.Println(betas[0][i] * alphas[0][i])

		newMu.Set(0, i, betas[0][i]*alphas[0][i]/divider)
	}

	for i := 0; i < states; i++ {
		for j := 0; j < states; j++ {
			bottom := 0.0
			top := 0.0

			for k := 0; k < states; k++ {
				tmp := 0.0
				for t := 0; t < T-1; t++ {
					tmp += alphas[t][i] * betas[t+1][k] * l.A.At(i, k) * l.CalcB(k, l.frames[t+1])
				}

				bottom += tmp
			}

			for t := 0; t < T-1; t++ {
				top += alphas[t][i] * betas[t+1][j] * l.A.At(i, j) * l.CalcB(j, l.frames[t+1])
			}

			newA.Set(i, j, top/bottom)
		}
	}

	for i := 0; i < states; i++ {
		newProcesses[i] = map[poisson.Region]float64{}
		// in formulas region is j
		for region := range *l.frames[0] {
			top := 0.0
			bottom := 0.0
			for t := 0; t < T; t++ {
				count := (*l.frames[t])[region]
				bottom += alphas[t][i] * betas[t][i]
				top += alphas[t][i] * betas[t][i] * float64(count)
			}

			newProcesses[i][region] = top / bottom
		}
	}

	l.Mu = newMu
	l.A = newA
	l.Processes = newProcesses
}
