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

	chainOnce, bdOnce, hdOnce, odOnce sync.Once
	rc                                []*ResultChain
	bd, hd                            []*mat.Dense
	od                                [][]poisson.IntensityMap
}

func NewRandomFactory(loadChain, loadDistribution bool) *RandomFactory {
	return &RandomFactory{
		LoadChain:        loadChain,
		LoadDistribution: loadDistribution,
	}
}

func (r *RandomFactory) Chains() []*ResultChain {
	r.chainOnce.Do(func() {
		if !r.LoadChain {
			for i := 0; i < constants.ObservationSequences; i++ {
				r.rc = append(r.rc, MarkovModel.Generate(constants.SampleSize))
			}

			utils.MustSave(constants.CachePath(constants.ChainFile), r.rc, json.Marshal)
		} else {
			utils.MustLoad(constants.CachePath(constants.ChainFile), &r.rc, json.Unmarshal)
		}
	})

	return lo.Map(r.rc, func(item *ResultChain, index int) *ResultChain {
		item.Frames = item.Frames[:constants.SampleSize]
		item.StateChain = item.StateChain[:constants.SampleSize]

		return item
	})
}

func (r *RandomFactory) BaseDistributions() []*mat.Dense {
	r.bdOnce.Do(func() {
		if !r.LoadDistribution {
			for i := 0; i < constants.LearnersCount; i++ {
				r.bd = append(r.bd, mat.NewDense(1, 2, utils.RandomStochasticMatrix(1, 2)))
			}

			utils.MustSave(constants.CachePath(constants.BaseDistributionFile), r.bd, utils.ManyDenseMarshal(r.bd))

		} else {
			utils.MustLoad(constants.CachePath(constants.BaseDistributionFile), &r.bd, utils.ManyDenseUnmarshal(&r.bd))
		}
	})

	return r.bd
}

func (r *RandomFactory) HiddenDistribution() []*mat.Dense {
	r.hdOnce.Do(func() {
		if !r.LoadDistribution {
			for i := 0; i < constants.LearnersCount; i++ {
				r.hd = append(r.hd, mat.NewDense(2, 2, utils.RandomStochasticMatrix(2, 2)))
			}

			utils.MustSave(constants.CachePath(constants.HiddenDistributionFile), r.hd, utils.ManyDenseMarshal(r.hd))
		} else {
			utils.MustLoad(constants.CachePath(constants.HiddenDistributionFile), &r.hd, utils.ManyDenseUnmarshal(&r.hd))
		}
	})

	return r.hd
}

func (r *RandomFactory) ObservableDistribution() [][]poisson.IntensityMap {
	return r.observableDistribution(utils.RandomMap[poisson.Region])
}

func (r *RandomFactory) ObservableUniformDistribution() [][]poisson.IntensityMap {
	return r.observableDistribution(utils.UniformMap[poisson.Region])
}

type RandomFunc func(keys []poisson.Region, expectedSum float64) map[poisson.Region]float64

func (r *RandomFactory) observableDistribution(randomFunc RandomFunc) [][]poisson.IntensityMap {
	r.odOnce.Do(func() {
		rc := r.Chains()
		totalPoints := lo.Reduce(rc[0].Frames, func(agg int, item *poisson.Area, index int) int {
			return item.TotalPoint() + agg
		}, 0)

		avgPoints := float64(totalPoints) / float64(len(rc[0].Frames))

		if !r.LoadDistribution {
			for i := 0; i < constants.LearnersCount; i++ {
				r.od = append(r.od, []poisson.IntensityMap{
					randomFunc(lo.Keys(P1), avgPoints),
					randomFunc(lo.Keys(P1), avgPoints),
				})
			}

			utils.MustSave(constants.CachePath(constants.ObservableDistributionFile), r.od, json.Marshal)
		} else {
			utils.MustLoad(constants.CachePath(constants.ObservableDistributionFile), &r.od, json.Unmarshal)
		}
	})

	return r.od
}
