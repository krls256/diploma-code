package markov

import (
	"fmt"
	"math"
	"testing"
)

var rc *ResultChain

const (
	T         = 10
	Threshold = 1e-10
)

func init() {
	rc = MarkovModel.Generate(T)
}

func TestForwardBackward(t *testing.T) {
	alphas, scalingCoeffs := Forward(false, rc.Frames, MarkovModel)
	betas := Backward(false, scalingCoeffs, rc.Frames, MarkovModel)

	_, stateCount := Mu.Dims()

	estimate1 := 0.0
	for i := 0; i < stateCount; i++ {
		estimate1 += MarkovModel.BaseDistribution.At(0, i) *
			PoissonProbability(i, rc.Frames[0], MarkovModel.ObservableProcesses) *
			betas[0][i]
	}

	estimate2 := 0.0
	for i := 0; i < stateCount; i++ {
		estimate2 += alphas[T-1][i]
	}

	estimate3 := 0.0
	for i := 0; i < stateCount; i++ {
		estimate3 += alphas[2][i] * betas[2][i]
	}

	// bad metric
	if math.Abs(estimate1-estimate2) > Threshold {
		t.Errorf("%v != %v", estimate1, estimate2)
	}

	if math.Abs(estimate2-estimate3) > Threshold {
		t.Errorf("%v != %v", estimate2, estimate3)
	}

	fmt.Println(estimate1)
	fmt.Println(estimate2)
	fmt.Println(estimate3)
}
