package main

import (
	"diploma/constants"
	"diploma/markov"
	"diploma/poisson"
	"diploma/utils"
	"encoding/json"
	"fmt"
	"sync"
)

func main() {
	rf := markov.NewRandomFactory(true, true)

	wg := sync.WaitGroup{}

	ms := markov.NewModelSelector(rf)

	rc := rf.Chains()
	wg.Add(len(rc))

	for i, item := range rc {
		go func(r *markov.ResultChain, i int) {
			poisson.DrawGif(r.Frames, r.StateChain,
				[]poisson.IntensityMap{markov.P1, markov.P2}, fmt.Sprintf("./report/chain-%v", i+1))

			wg.Done()
		}(item, i)
	}

	ls := ms.Run()

	for j, intensity := range markov.MarkovModel.ObservableProcesses {
		path := constants.CachePath(fmt.Sprintf("org-"+constants.IntensityMapFormat, 0, j))

		utils.MustSave(path, intensity, json.Marshal)
	}

	for i, l := range ls {
		l.StartModel.Draw(fmt.Sprintf("report/generated-%v.png", i+1))
		l.Model.Draw(fmt.Sprintf("report/learned-%v.png", i+1))
		l.Model.SortByStandardDeviation(markov.MarkovModel)
		l.Model.Draw(fmt.Sprintf("report/learned-sort-%v.png", i+1))

		utils.DrawValuesGraphics(l.LogProbs, "log(P(Y = y))", fmt.Sprintf("log-prob-%v", i+1))

		for j, intensity := range l.Model.ObservableProcesses {
			path := constants.CachePath(fmt.Sprintf(constants.IntensityMapFormat, i, j))

			utils.MustSave(path, intensity, json.Marshal)
		}
	}

	wg.Wait()

	markov.MarkovModel.Draw("report/original.png")
}
