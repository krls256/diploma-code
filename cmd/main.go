package main

import (
	"diploma/markov"
	"diploma/poisson"
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"sync"
)

const (
	N = 100

	LearnerCount = 1
	Iters        = 100
)

func main() {
	rf := markov.NewRandomFactory(false, false)
	rc := rf.Chain()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		poisson.DrawGif(rc.Frames, rc.StateChain, []poisson.IntensityMap{markov.P1, markov.P2}, "./report/chain")

		wg.Done()
	}()

	for i := 0; i < LearnerCount; i++ {
		wg.Add(1)

		go func(index int) {
			before, after := learn(rc.Frames, rf)

			before.Draw(fmt.Sprintf("report/generated-%v.png", index+1))
			after.Draw(fmt.Sprintf("report/learned-%v.png", index+1))

			wg.Done()
		}(i)
	}

	wg.Wait()

	markov.MarkovModel.Draw("report/original.png")
}

func learn(frames []*poisson.Area, rf *markov.RandomFactory) (before, after *markov.Model) {
	muL := rf.BaseDistribution()
	aL := rf.HiddenDistribution()
	processL := rf.ObservableDistribution()

	l := markov.NewLearner(markov.NewModel(muL, aL, processL), frames)

	before = l.Model.DeepCopy()

	logProbs := []float64{}

	for i := 0; i < Iters; i++ {
		l.Step()
		logProbs = append(logProbs, l.LogProb())
	}

	drawLogProb(logProbs)

	return before, l.Model
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
