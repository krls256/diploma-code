package poisson

import (
	"bytes"
	"diploma/utils"
	"github.com/samber/lo"
	"math/rand"
	"sort"
)

const (
	StringCellWidth = 8
)

type ProcessWithIntensityFunc struct {
	XStart, XEnd, YStart, YEnd float64
	IntensityFunc              func(x, y float64) float64
	XAxisPartQuantity          int
	YAxisPartQuantity          int
}

type IntensityMap map[Region]float64

func (im IntensityMap) String() string {
	buf := bytes.Buffer{}

	regions := lo.Keys(im)
	sort.Slice(regions, func(i, j int) bool {
		return RegionRowOrderBiggerThan(regions[i], regions[j])
	})

	for i := 0; i < len(regions); i++ {
		if i > 0 && regions[i-1][3] != regions[i][3] {
			buf.WriteString("\n")
		}
		str := utils.PrecisionString(im[regions[i]])

		buf.WriteString(utils.WrapInBox(str, StringCellWidth))
	}

	return buf.String()
}

func (im IntensityMap) Copy() IntensityMap {
	newIm := IntensityMap{}

	for k, v := range im {
		newIm[k] = v
	}

	return newIm
}

func (im IntensityMap) Generate() *Area {
	area := Area{}

	intensity := 0.0

	for _, in := range im {
		intensity += in
	}

	process := poisson(intensity)
	points := int(process.Rand())

	for i := 0; i < points; i++ {
		r := rand.Float64() * intensity

		for reg, in := range im {
			if _, ok := area[reg]; !ok {
				area[reg] = 0
			}

			if r < in {
				area[reg]++

				break
			}

			r -= in
		}
	}

	return &area
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

func (p *ProcessWithIntensityFunc) IntensityMap() IntensityMap {
	xStep := (p.XEnd - p.XStart) / float64(p.XAxisPartQuantity)
	yStep := (p.YEnd - p.YStart) / float64(p.YAxisPartQuantity)

	intensityMap := IntensityMap{}

	for xStart, xEnd := p.XStart, p.XStart+xStep; xEnd <= p.XEnd; xStart, xEnd = xStart+xStep, xEnd+xStep {
		for yStart, yEnd := p.YStart, p.YStart+yStep; yEnd <= p.YEnd; yStart, yEnd = yStart+yStep, yEnd+yStep {

			middleX, middleY := xStart+(xEnd-xStart)/2, yStart+(yEnd-yStart)/2
			intensityMap[Region{xStart, xEnd, yStart, yEnd}] = p.IntensityFunc(middleX, middleY)
		}
	}

	return intensityMap
}

func (p *ProcessWithIntensityFunc) Generate() *Area {
	return p.IntensityMap().Generate()
}
