package markov

import (
	"diploma/constants"
	"diploma/poisson"
	"diploma/utils"
	"encoding/json"
	"github.com/samber/lo"
	"gonum.org/v1/gonum/mat"
	"sync"
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

func (r *RandomFactory) Chain() *ResultChain {
	r.chainOnce.Do(func() {
		if !r.LoadChain {
			r.rc = MarkovModel.Generate(constants.SampleSize)
			utils.MustSave(constants.CachePath(constants.ChainFile), r.rc, json.Marshal)
		} else {
			utils.MustLoad(constants.CachePath(constants.ChainFile), r.rc, json.Unmarshal)
		}
	})

	return r.rc
}

func (r *RandomFactory) BaseDistribution() *mat.Dense {
	baseDistribution := &mat.Dense{}

	if !r.LoadDistribution {
		baseDistribution = mat.NewDense(1, 2, utils.RandomStochasticMatrix(1, 2))
		utils.MustSave(constants.CachePath(constants.BaseDistributionFile), baseDistribution, utils.DenseMarshal(baseDistribution))

	} else {
		utils.MustLoad(constants.CachePath(constants.BaseDistributionFile), baseDistribution, utils.DenseUnmarshal(baseDistribution))
	}

	return baseDistribution
}

func (r *RandomFactory) HiddenDistribution() *mat.Dense {
	hiddenDistribution := &mat.Dense{}

	if !r.LoadDistribution {
		hiddenDistribution = mat.NewDense(2, 2, utils.RandomStochasticMatrix(2, 2))
		utils.MustSave(constants.CachePath(constants.HiddenDistributionFile), hiddenDistribution, utils.DenseMarshal(hiddenDistribution))
	} else {
		hiddenDistribution = &mat.Dense{}
		utils.MustLoad(constants.CachePath(constants.HiddenDistributionFile), hiddenDistribution, utils.DenseUnmarshal(hiddenDistribution))
	}

	return hiddenDistribution
}

func (r *RandomFactory) ObservableDistribution() []poisson.IntensityMap {
	return r.observableDistribution(utils.RandomMap[poisson.Region])
}

func (r *RandomFactory) ObservableUniformDistribution() []poisson.IntensityMap {
	return r.observableDistribution(utils.UniformMap[poisson.Region])
}

type RandomFunc func(keys []poisson.Region, expectedSum float64) map[poisson.Region]float64

func (r *RandomFactory) observableDistribution(randomFunc RandomFunc) []poisson.IntensityMap {
	observableDistribution := []poisson.IntensityMap{}

	rc := r.Chain()
	totalPoints := lo.Reduce(rc.Frames, func(agg int, item *poisson.Area, index int) int {
		return item.TotalPoint() + agg
	}, 0)

	avgPoints := float64(totalPoints) / float64(len(rc.Frames))

	if !r.LoadDistribution {
		observableDistribution = append(observableDistribution,
			randomFunc(lo.Keys(P1), avgPoints),
			randomFunc(lo.Keys(P1), avgPoints))

		utils.MustSave(constants.CachePath(constants.ObservableDistributionFile), observableDistribution, json.Marshal)
	} else {
		utils.MustLoad(constants.CachePath(constants.ObservableDistributionFile), &observableDistribution, json.Unmarshal)
	}

	return observableDistribution
}
