package main

import (
	"diploma/markov"
	"diploma/poisson"
	"diploma/utils"
	"fmt"
	"gonum.org/v1/gonum/mat"
	"math"
)

const (
	xStartG           = 0
	xEndG             = 10
	yStartG           = 0
	yEndG             = 10
	XAxisPartQuantity = 2
	YAxisPartQuantity = 2

	xCenter = (xEndG - xStartG) / 2
	yCenter = (yEndG - yStartG) / 2

	N = 50
)

func regular(x, y float64) float64 {
	return math.Max(-math.Pow(x-y, 2)/2+30, 0)
}

func irregular(x, y float64) float64 {
	return math.Max(math.Pow(y-yCenter, 3)+10, 4)
}

func cand1(x, y float64) float64 {
	return y * y / 10
}

func cand2(x, y float64) float64 {
	return math.Pow(math.E, 4+(y-x)*(x-y)/24) / 15
}

var (
	p1 = poisson.NewProcessWithIntensityFunc(xStartG, xEndG, yStartG, yEndG, XAxisPartQuantity, YAxisPartQuantity, cand1)
	p2 = poisson.NewProcessWithIntensityFunc(xStartG, xEndG, yStartG, yEndG, XAxisPartQuantity, YAxisPartQuantity, cand2)
	mu = mat.NewDense(1, 2, []float64{0.75, 0.25})
	a  = mat.NewDense(2, 2, []float64{
		0.65, 0.35,
		0.5, 0.5,
	})

	proc = markov.NewModel(mu, a, []*poisson.ProcessWithIntensityFunc{p1, p2})
)

func main() {
	frames, statesChain := proc.Generate(N)

	muL := mat.NewDense(1, 2, []float64{0.70, 0.30})
	aL := mat.NewDense(2, 2, []float64{
		0.65, 0.35,
		0.45, 0.55,
	})

	process := []map[poisson.Region]float64{
		utils.RandomizeMap(p2.IntensityMap(), 2),
		utils.RandomizeMap(p2.IntensityMap(), 3)}
	fmt.Println("1 mmm", p1.IntensityMap())
	fmt.Println("2 mmm", p2.IntensityMap())
	l := markov.NewLearner(muL, aL, process, frames)

	for i := 0; i < 200; i++ {
		//fmt.Println("---------------- new step")
		l.Step()
	}

	fmt.Println("mu:", l.Mu)
	fmt.Println("a", l.A)
	fmt.Println("process 1", l.Processes[0])
	fmt.Println("process 2", l.Processes[1])

	poisson.DrawGif(frames, statesChain, []map[poisson.Region]float64{p1.IntensityMap(), p2.IntensityMap()}, "./tmp/card")
}

func formatNum(num float64) string {
	if num > 100 {
		return fmt.Sprintf("%3.0f", num)
	}

	if num > 10 {
		return fmt.Sprintf("%3.1f", num)
	}

	return fmt.Sprintf("%3.2f", num)
}
