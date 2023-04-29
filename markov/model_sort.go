package markov

import "diploma/poisson"

type ModelSorter struct {
	toSort     *Model
	baseOnSort *Model
}

func (ms *ModelSorter) Len() int {
	return len(ms.toSort.ObservableProcesses)
}

func (ms *ModelSorter) Less(i, j int) bool {
	a := poisson.StandardDeviation(ms.toSort.ObservableProcesses[i], ms.baseOnSort.ObservableProcesses[i])
	b := poisson.StandardDeviation(ms.toSort.ObservableProcesses[j], ms.baseOnSort.ObservableProcesses[i])
	return a > b
}

func (ms *ModelSorter) Swap(i, j int) {
	a, b := ms.toSort.BaseDistribution.At(0, i), ms.toSort.BaseDistribution.At(0, j)

	ms.toSort.BaseDistribution.Set(0, i, b)
	ms.toSort.BaseDistribution.Set(0, j, a)

	rows, cols := ms.toSort.HiddenDistribution.Dims()
	for k := 0; k < rows; k++ {
		a, b = ms.toSort.HiddenDistribution.At(k, i), ms.toSort.HiddenDistribution.At(k, j)

		ms.toSort.HiddenDistribution.Set(k, i, b)
		ms.toSort.HiddenDistribution.Set(k, j, a)
	}

	for k := 0; k < cols; k++ {
		a, b = ms.toSort.HiddenDistribution.At(i, k), ms.toSort.HiddenDistribution.At(j, k)

		ms.toSort.HiddenDistribution.Set(i, k, b)
		ms.toSort.HiddenDistribution.Set(j, k, a)
	}

	ms.toSort.ObservableProcesses[i], ms.toSort.ObservableProcesses[j] = ms.toSort.ObservableProcesses[j], ms.toSort.ObservableProcesses[i]
}
