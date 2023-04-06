package poisson

import (
	src "golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func poisson(intensity float64) *distuv.Poisson {
	return &distuv.Poisson{Lambda: intensity, Src: src.NewSource(rand.Uint64())}
}
