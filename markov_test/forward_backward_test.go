package markov_test

import (
	"diploma/config"
	"diploma/markov"
	"fmt"
	"math"
	"testing"
)

var rc *markov.ResultChain

const (
	T         = 10
	Threshold = 1e-10
)

func init() {
	rc = config.MarkovModel.Generate(T)
}

func TestForwardBackward(t *testing.T) {
	alphas, scalingCoeffs := markov.Forward(false, rc.Frames, config.MarkovModel)
	betas := markov.Backward(false, scalingCoeffs, rc.Frames, config.MarkovModel)

	_, stateCount := config.Mu.Dims()

	estimate1 := 0.0
	for i := 0; i < stateCount; i++ {
		estimate1 += config.MarkovModel.Mu.At(0, i) *
			markov.PoissonProbability(i, rc.Frames[0], config.MarkovModel.HiddenProcesses) *
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
