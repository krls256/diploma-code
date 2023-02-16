package poisson

import (
	"math/rand"
)

type ProcessWithIntensityFunc struct {
	XStart, XEnd, YStart, YEnd float64
	IntensityFunc              func(x, y float64) float64
	XAxisPartQuantity          int
	YAxisPartQuantity          int
}

func NewProcessWithIntensityFunc(
	XStart, XEnd, YStart, YEnd float64,
	XAxisPartQuantity int, YAxisPartQuantity int,
	IntensityFunc func(x, y float64) float64) *ProcessWithIntensityFunc {
	return &ProcessWithIntensityFunc{
		XStart: XStart, XEnd: XEnd, YStart: YStart, YEnd: YEnd,
		XAxisPartQuantity: XAxisPartQuantity,
		YAxisPartQuantity: YAxisPartQuantity,
		IntensityFunc:     IntensityFunc,
	}
}

func (p *ProcessWithIntensityFunc) IntensityMap() map[Region]float64 {
	xStep := (p.XEnd - p.XStart) / float64(p.XAxisPartQuantity)
	yStep := (p.YEnd - p.YStart) / float64(p.YAxisPartQuantity)

	intensityMap := map[Region]float64{}

	for xStart, xEnd := p.XStart, p.XStart+xStep; xEnd <= p.XEnd; xStart, xEnd = xStart+xStep, xEnd+xStep {
		for yStart, yEnd := p.YStart, p.YStart+yStep; yEnd <= p.YEnd; yStart, yEnd = yStart+yStep, yEnd+yStep {

			middleX, middleY := xStart+(xEnd-xStart)/2, yStart+(yEnd-yStart)/2
			intensityMap[Region{xStart, xEnd, yStart, yEnd}] = p.IntensityFunc(middleX, middleY)
		}
	}

	return intensityMap
}

func (p *ProcessWithIntensityFunc) Generate() *Area {
	xStep := (p.XEnd - p.XStart) / float64(p.XAxisPartQuantity)
	yStep := (p.YEnd - p.YStart) / float64(p.YAxisPartQuantity)

	area := Area{}
	intensityMap := map[Region]float64{}

	for xStart, xEnd := p.XStart, p.XStart+xStep; xEnd <= p.XEnd; xStart, xEnd = xStart+xStep, xEnd+xStep {
		for yStart, yEnd := p.YStart, p.YStart+yStep; yEnd <= p.YEnd; yStart, yEnd = yStart+yStep, yEnd+yStep {

			middleX, middleY := xStart+(xEnd-xStart)/2, yStart+(yEnd-yStart)/2
			intensityMap[Region{xStart, xEnd, yStart, yEnd}] = p.IntensityFunc(middleX, middleY)
			area[Region{xStart, xEnd, yStart, yEnd}] = 0
		}
	}

	intensity := 0.0

	for _, in := range intensityMap {
		intensity += in
	}

	process := poisson(intensity)
	points := int(process.Rand())

	for i := 0; i < points; i++ {
		r := rand.Float64() * intensity

		for reg, in := range intensityMap {
			if r < in {
				area[reg]++

				break
			}

			r -= in
		}
	}

	return &area
}
