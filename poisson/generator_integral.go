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

type ProcessWithIntegralIntensity struct {
	XStart, XEnd, YStart, YEnd float64
	IntensityIntegral          func(xStart, xEnd, yStart, yEnd float64) float64
	XAxisPartQuantity          int
	YAxisPartQuantity          int
}

func NewProcessWithIntegralIntensity(
	XStart, XEnd, YStart, YEnd float64,
	XAxisPartQuantity int, YAxisPartQuantity int,
	IntensityIntegral func(xStart, xEnd, yStart, yEnd float64) float64) *ProcessWithIntegralIntensity {
	return &ProcessWithIntegralIntensity{
		XStart: XStart, XEnd: XEnd, YStart: YStart, YEnd: YEnd,
		XAxisPartQuantity: XAxisPartQuantity,
		YAxisPartQuantity: YAxisPartQuantity,
		IntensityIntegral: IntensityIntegral,
	}
}

func (p *ProcessWithIntegralIntensity) Generate() *Area {
	xStep := (p.XEnd - p.XStart) / float64(p.XAxisPartQuantity)
	yStep := (p.YEnd - p.YStart) / float64(p.YAxisPartQuantity)

	area := Area{}

	for xStart, xEnd := p.XStart, p.XStart+xStep; xEnd <= p.XEnd; xStart, xEnd = xStart+xStep, xEnd+xStep {
		for yStart, yEnd := p.YStart, p.YStart+yStep; yEnd <= p.YEnd; yStart, yEnd = yStart+yStep, yEnd+yStep {

			intensity := p.IntensityIntegral(xStart, xEnd, yStart, yEnd)
			process := poisson(intensity)
			area[Region{xStart, xEnd, yStart, yEnd}] = int(process.Rand())
		}
	}

	return &area
}

func poisson(intensity float64) *distuv.Poisson {
	return &distuv.Poisson{Lambda: intensity, Src: src.NewSource(rand.Uint64())}
}
