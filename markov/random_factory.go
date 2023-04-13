package markov

import (
	"diploma/poisson"
	"diploma/utils"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"gonum.org/v1/gonum/mat"
	"sync"
)

const (
	CacheDir                   = "./cache"
	ChainFile                  = "chain.json"
	BaseDistributionFile       = "base-distribution.bin"
	HiddenDistributionFile     = "hidden-distribution.bin"
	ObservableDistributionFile = "observable-distribution.json"
)

type RandomFactory struct {
	LoadChain        bool
	LoadDistribution bool

	chainOnce sync.Once
	rc        *ResultChain
}

func NewRandomFactory(loadChain, loadDistribution bool) *RandomFactory {
	return &RandomFactory{
		LoadChain:        loadChain,
		LoadDistribution: loadDistribution,
	}
}

func (r *RandomFactory) path(filename string) string {
	return fmt.Sprintf("%s/%s", CacheDir, filename)
}

func (r *RandomFactory) Chain() *ResultChain {
	r.chainOnce.Do(func() {
		if !r.LoadChain {
			r.rc = MarkovModel.Generate(SampleSize)
			utils.MustSave(r.path(ChainFile), r.rc, json.Marshal)
		} else {
			utils.MustLoad(r.path(ChainFile), r.rc, json.Unmarshal)
		}
	})

	return r.rc
}

func (r *RandomFactory) BaseDistribution() *mat.Dense {
	baseDistribution := &mat.Dense{}

	if !r.LoadDistribution {
		baseDistribution = mat.NewDense(1, 2, utils.RandomStochasticMatrix(1, 2))
		utils.MustSave(r.path(BaseDistributionFile), baseDistribution, utils.DenseMarshal(baseDistribution))

	} else {
		utils.MustLoad(r.path(BaseDistributionFile), baseDistribution, utils.DenseUnmarshal(baseDistribution))
	}

	return baseDistribution
}

func (r *RandomFactory) HiddenDistribution() *mat.Dense {
	hiddenDistribution := &mat.Dense{}

	if !r.LoadDistribution {
		hiddenDistribution = mat.NewDense(2, 2, utils.RandomStochasticMatrix(2, 2))
		utils.MustSave(r.path(HiddenDistributionFile), hiddenDistribution, utils.DenseMarshal(hiddenDistribution))
	} else {
		hiddenDistribution = &mat.Dense{}
		utils.MustLoad(r.path(HiddenDistributionFile), hiddenDistribution, utils.DenseUnmarshal(hiddenDistribution))
	}

	return hiddenDistribution
}

func (r *RandomFactory) ObservableDistribution() []poisson.IntensityMap {
	observableDistribution := []poisson.IntensityMap{}

	rc := r.Chain()
	totalPoints := lo.Reduce(rc.Frames, func(agg int, item *poisson.Area, index int) int {
		return item.TotalPoint() + agg
	}, 0)

	avgPoints := float64(totalPoints) / float64(len(rc.Frames))

	if !r.LoadDistribution {
		observableDistribution = append(observableDistribution,
			utils.RandomMap(lo.Keys(P1), avgPoints),
			utils.RandomMap(lo.Keys(P1), avgPoints))

		utils.MustSave(r.path(ObservableDistributionFile), observableDistribution, json.Marshal)
	} else {
		utils.MustLoad(r.path(ObservableDistributionFile), &observableDistribution, json.Unmarshal)
	}

	return observableDistribution
}
