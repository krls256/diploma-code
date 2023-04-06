package main

import (
	"diploma/config"
	"diploma/markov"
	"diploma/poisson"
	"diploma/utils"
	"encoding/json"
	"fmt"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"sync"
)

const (
	N = 100

	LearnerCount = 1
	Iters        = 500
)

func main() {
	//rc := config.MarkovModel.Generate(N)
	//utils.MustSave("./cache/data.json", rc, json.Marshal)

	rc := markov.ResultChain{}
	utils.MustLoad("./cache/data.json", &rc, json.Unmarshal)

	//go poisson.DrawGif(rc.Frames, rc.StateChain, []poisson.IntensityMap{p1, p2}, "./tmp/card")

	wg := sync.WaitGroup{}

	for i := 0; i < LearnerCount; i++ {
		wg.Add(1)

		go learn(rc.Frames, &wg)
	}

	wg.Wait()

	//fmt.Println("Waiting 10 seconds for gif")
	//time.Sleep(time.Second * 10)
}

func learn(frames []*poisson.Area, wg *sync.WaitGroup) {
	muL := mat.NewDense(1, 2, utils.RandomStochasticMatrix(1, 2))
	utils.MustSave("./cache/mu.bin", muL, utils.DenseMarshal(muL))

	//muL := &mat.Dense{}
	//utils.MustLoad("./cache/mu.bin", muL, utils.DenseUnmarshal(muL))

	aL := mat.NewDense(2, 2, utils.RandomStochasticMatrix(2, 2))
	utils.MustSave("./cache/a.bin", aL, utils.DenseMarshal(aL))

	//aL := &mat.Dense{}
	//utils.MustLoad("./cache/a.bin", aL, utils.DenseUnmarshal(aL))

	processL := []poisson.IntensityMap{
		config.P1.Copy(),
		config.P2.Copy(),
		//utils.RandomMap(lo.Keys(config.P1.Copy()), 5, 5),
		//utils.RandomMap(lo.Keys(config.P2.Copy()), 5, 5),
	}
	utils.MustSave("./cache/process.json", processL, json.Marshal)
	//utils.MustLoad("./cache/process.json", &processL, json.Unmarshal)

	l := markov.NewLearner(markov.NewModel(muL, aL, processL), frames)

	//fmt.Println(l.Model.String())

	logProbs := []float64{}

	for i := 0; i < Iters; i++ {
		l.Step()
		logProbs = append(logProbs, l.LogProb())
		fmt.Println(l.Model.String())
	}

	drawLogProb(logProbs)

	wg.Done()
}

func drawLogProb(probs []float64) {
	p := plot.New()

	p.Title.Text = "Log Prob"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(plotter.NewGrid())

	lineData := plotter.XYs{}
	for i, pr := range probs {
		lineData = append(lineData, plotter.XY{X: float64(i), Y: pr})
	}

	l, err := plotter.NewLine(lineData)
	if err != nil {
		panic(err)
	}
	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	l.LineStyle.Color = color.RGBA{B: 255, A: 255}

	p.Add(l)

	if err = p.Save(4*vg.Inch, 4*vg.Inch, "log-prob.png"); err != nil {
		panic(err)
	}
}
