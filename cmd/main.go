package main

import (
	"diploma/markov"
	"diploma/poisson"
	"diploma/utils"
	"fmt"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/mat"
	"math"

	"time"
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

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}

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
	p1 = poisson.NewProcessWithIntensityFunc(xStartG, xEndG, yStartG, yEndG, XAxisPartQuantity, YAxisPartQuantity, cand1).IntensityMap()
	p2 = poisson.NewProcessWithIntensityFunc(xStartG, xEndG, yStartG, yEndG, XAxisPartQuantity, YAxisPartQuantity, cand2).IntensityMap()
	mu = mat.NewDense(1, 2, []float64{0.75, 0.25})
	a  = mat.NewDense(2, 2, []float64{
		0.65, 0.35,
		0.5, 0.5,
	})

	proc = markov.NewModel(mu, a, []poisson.IntensityMap{p1, p2})
)

func main() {
	frames, statesChain := proc.Generate(N)

	muL := mat.NewDense(1, 2, []float64{0.70, 0.30})
	aL := mat.NewDense(2, 2, []float64{
		0.65, 0.35,
		0.45, 0.55,
	})

	process := []poisson.IntensityMap{
		utils.RandomizeMap(p2.Copy(), 1),
		utils.RandomizeMap(p2.Copy(), 1)}

	l := markov.NewLearner(markov.NewModel(muL, aL, process), frames)

	for i := 0; i < 20; i++ {
		l.Step()
	}

	fmt.Println(l.Model)

	poisson.DrawGif(frames, statesChain, []poisson.IntensityMap{p1, p2}, "./tmp/card")
}
