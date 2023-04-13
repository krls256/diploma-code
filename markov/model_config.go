package markov

import (
	"diploma/poisson"
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
	XAxisPartQuantity = 3
	YAxisPartQuantity = 3

	xCenter = (xEndG - xStartG) / 2
	yCenter = (yEndG - yStartG) / 2

	SampleSize = 100
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
	P1 = poisson.NewProcessWithIntensityFunc(xStartG, xEndG, yStartG, yEndG, XAxisPartQuantity, YAxisPartQuantity, cand1).IntensityMap()
	P2 = poisson.NewProcessWithIntensityFunc(xStartG, xEndG, yStartG, yEndG, XAxisPartQuantity, YAxisPartQuantity, cand2).IntensityMap()
	Mu = mat.NewDense(1, 2, []float64{0.5, 0.5})
	A  = mat.NewDense(2, 2, []float64{
		0.65, 0.35,
		0.5, 0.5,
	})

	MarkovModel = NewModel(Mu, A, []poisson.IntensityMap{P1, P2})
)
