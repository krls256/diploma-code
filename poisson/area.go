package poisson

import (
	"fmt"
	"github.com/samber/lo"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"math/rand"
	"time"
)

type Area map[Region]int
type Region [4]float64 // XStart, XEnd, YStart, YEnd

func (r Region) String() string {
	return fmt.Sprintf("%v-%v-%v-%v", r[0], r[1], r[2], r[3])
}

func (r Region) Center() (xCenter, yCenter float64) {
	xStart, xEnd, yStart, yEnd := r[0], r[1], r[2], r[3]
	return xStart + (xEnd-xStart)/2, yStart + (yEnd-yStart)/2
}

func (a *Area) TotalPoint() int {
	return lo.Reduce(lo.Values(*a), func(agg int, item int, index int) int {
		return agg + item
	}, 0)
}

func (a *Area) Draw(filename string) {
	rand.Seed(time.Now().Unix())
	scatterData := a.Points()

	p := DrawGrid(lo.Keys(*a), "Points Process")

	s, err := plotter.NewScatter(scatterData)
	if err != nil {
		panic(err)
	}
	s.Radius = vg.Points(1)
	s.GlyphStyle.Color = color.RGBA{A: 255}

	p.Legend.Add("scatter", s)

	p.Add(s)

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, filename); err != nil {
		panic(err)
	}
}

func (a *Area) Points() plotter.XYs {
	pts := make(plotter.XYs, 0)

	for region, count := range *a {
		for i := 0; i < count; i++ {
			x, y := randomXY(region)
			p := plotter.XY{X: x, Y: y}

			pts = append(pts, p)
		}
	}

	return pts
}

func (a *Area) RegionLines() []*plotter.Line {
	lines := make([]*plotter.Line, 0, len(*a))
	for region := range *a {
		XStart, XEnd, YStart, YEnd := region[0], region[1], region[2], region[3]
		l, err := plotter.NewLine(plotter.XYs{
			plotter.XY{X: XStart, Y: YStart},
			plotter.XY{X: XStart, Y: YEnd},
			plotter.XY{X: XEnd, Y: YEnd},
			plotter.XY{X: XEnd, Y: YStart},
		})

		if err != nil {
			panic(err)
		}

		l.LineStyle.Width = vg.Points(1)
		l.LineStyle.Color = color.RGBA{G: 255, A: 255}
		lines = append(lines, l)

	}

	return lines
}

func randomXY(region Region) (x, y float64) {
	XStart, XEnd, YStart, YEnd := region[0], region[1], region[2], region[3]
	xSize, ySize := XEnd-XStart, YEnd-YStart
	relX, relY := rand.Float64()*xSize, rand.Float64()*ySize

	return relX + XStart, relY + YStart
}

func DrawGrid(regions []Region, title string) *plot.Plot {
	p := plot.New()

	p.Title.Text = title
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(plotter.NewGrid())

	grid := RegionLines(regions)
	for _, l := range grid {
		p.Add(l)
	}

	return p
}

func RegionLines(regions []Region) []*plotter.Line {
	lines := make([]*plotter.Line, 0, len(regions))
	for _, region := range regions {
		XStart, XEnd, YStart, YEnd := region[0], region[1], region[2], region[3]
		l, err := plotter.NewLine(plotter.XYs{
			plotter.XY{X: XStart, Y: YStart},
			plotter.XY{X: XStart, Y: YEnd},
			plotter.XY{X: XEnd, Y: YEnd},
			plotter.XY{X: XEnd, Y: YStart},
		})

		if err != nil {
			panic(err)
		}

		l.LineStyle.Width = vg.Points(1)
		l.LineStyle.Color = color.RGBA{G: 255, A: 255}
		lines = append(lines, l)

	}

	return lines
}
