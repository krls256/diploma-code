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
	rf := markov.NewRandomFactory(false, false)

	wg := sync.WaitGroup{}

	ms := markov.NewModelSelector(rf)

	ls := ms.Run()

	rc := rf.Chain()
	wg.Add(1)

	go func() {
		poisson.DrawGif(rc.Frames, rc.StateChain,
			[]poisson.IntensityMap{markov.P1, markov.P2}, "./report/chain")

		wg.Done()
	}()

	for i, l := range ls {
		l.StartModel.Draw(fmt.Sprintf("report/generated-%v.png", i+1))
		l.Model.Draw(fmt.Sprintf("report/learned-%v.png", i+1))
		l.Model.SortByStandardDeviation(markov.MarkovModel)
		l.Model.Draw(fmt.Sprintf("report/learned-sort-%v.png", i+1))

		for j, intensity := range l.Model.ObservableProcesses {
			path := constants.CachePath(fmt.Sprintf(constants.IntensityMapFormat, i, j))

			utils.MustSave(path, intensity, json.Marshal)
		}
	}

	wg.Wait()

	markov.MarkovModel.Draw("report/original.png")
}
