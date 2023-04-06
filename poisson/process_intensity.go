package poisson

const (
	StringCellWidth = 8
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
