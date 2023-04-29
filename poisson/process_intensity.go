package poisson

import (
	"diploma/constants"
)

const (
	StringCellWidth = 8
)

type ProcessWithIntensityFunc struct {
	XStart, XEnd, YStart, YEnd float64
	IntensityFunc              func(x, y float64) float64
}

func NewProcessWithIntensityFunc(
	XStart, XEnd, YStart, YEnd float64,
	IntensityFunc func(x, y float64) float64) *ProcessWithIntensityFunc {
	return &ProcessWithIntensityFunc{
		XStart: XStart, XEnd: XEnd, YStart: YStart, YEnd: YEnd,
		IntensityFunc: IntensityFunc,
	}
}

func (p *ProcessWithIntensityFunc) IntensityMap() IntensityMap {

	var (
		im IntensityMap
		ok bool
	)

	for {
		im, ok = p.intensityMap()

		if !ok {
			constants.AddSplit()
		} else {
			break
		}
	}

	return im
}

func (p *ProcessWithIntensityFunc) intensityMap() (IntensityMap, bool) {
	xStep := (p.XEnd - p.XStart) / float64(constants.XAxisSplit)
	yStep := (p.YEnd - p.YStart) / float64(constants.YAxisSplit)

	intensityMap := IntensityMap{}

	for xStart, xEnd := p.XStart, p.XStart+xStep; xEnd <= p.XEnd; xStart, xEnd = xStart+xStep, xEnd+xStep {
		for yStart, yEnd := p.YStart, p.YStart+yStep; yEnd <= p.YEnd; yStart, yEnd = yStart+yStep, yEnd+yStep {

			//area := (xEnd - xStart) * (yEnd - yStart)
			//middleX, middleY := xStart+(xEnd-xStart)/2, yStart+(yEnd-yStart)/2
			//intensity := p.IntensityFunc(middleX, middleY) * area

			// square method
			intensity := 0.0
			area := constants.IntensityStep * constants.IntensityStep
			for xss := xStart; xss < xEnd; xss += constants.IntensityStep {
				for yss := yStart; yss < yEnd; yss += constants.IntensityStep {
					xssCenter, yssCenter := xss+constants.IntensityStep/2, yss+constants.IntensityStep/2

					intensity += p.IntensityFunc(xssCenter, yssCenter) * area
				}
			}

			if intensity > constants.MaxPoissonIntensity {
				return nil, false
			}

			intensityMap[Region{xStart, xEnd, yStart, yEnd}] = intensity
		}
	}

	return intensityMap, true
}

func (p *ProcessWithIntensityFunc) Generate() *Area {
	return p.IntensityMap().Generate()
}
