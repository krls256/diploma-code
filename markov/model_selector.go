package markov

import (
	"diploma/constants"
	"diploma/poisson"
	"fmt"
	"github.com/samber/lo"
	"github.com/schollz/progressbar/v3"
	"sort"
	"sync"
)

type ModelSelector struct {
	rf *RandomFactory
}

func NewModelSelector(rf *RandomFactory) *ModelSelector {
	return &ModelSelector{
		rf: rf,
	}
}

func (ms *ModelSelector) Run() []*Learner {
	wg := sync.WaitGroup{}

	learners := make([]*Learner, 0)

	for j := 0; j < constants.LearnersCount; j++ {
		learners = append(learners, ms.randomLearner())
	}

	bar := progressbar.Default(constants.LearnersCount * constants.ItersPerEpoch)

	frames := lo.Map(ms.rf.Chains(), func(item *ResultChain, index int) []*poisson.Area {
		return item.Frames
	})

	for j := 0; j < len(learners); j++ {
		wg.Add(1)

		learners[j].ObservationSequences = frames

		go func(index int) {
			ms.learn(learners[index], bar)

			wg.Done()
		}(j)
	}

	wg.Wait()

	sort.Slice(learners, func(i, j int) bool {
		return learners[i].LogProb() > learners[j].LogProb()
	})

	fmt.Println("Learning is finished")

	return learners
}

func (ms *ModelSelector) randomLearner() *Learner {
	muL := ms.rf.BaseDistribution()
	aL := ms.rf.HiddenDistribution()
	processL := ms.rf.ObservableDistribution()

	return NewLearner(NewModel(muL, aL, processL), nil)
}

func (ms *ModelSelector) learn(l *Learner, bar *progressbar.ProgressBar) {
	for i := 0; i < constants.ItersPerEpoch; i++ {
		l.Step()
		l.LogProbs = append(l.LogProbs, l.LogProb())

		if err := bar.Add(1); err != nil {
			panic(err)
		}
	}
}
