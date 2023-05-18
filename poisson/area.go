package poisson

import (
	"bytes"
	"diploma/utils"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image"
	"image/color"
	"strconv"
)

type Area map[Region]int

func (a *Area) TotalPoint() int {
	return lo.Reduce(lo.Values(*a), func(agg int, item int, index int) int {
		return agg + item
	}, 0)
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

func (a *Area) PointsImage(title string) image.Image {
	scatterData := a.Points()

	p := DrawGrid(lo.Keys(*a), title)

	s, err := plotter.NewScatter(scatterData)
	if err != nil {
		panic(err)
	}
	s.Radius = vg.Points(1)
	s.GlyphStyle.Color = color.RGBA{A: 255}

	p.Add(s)

	return utils.PlotToPNG(p)
}

func (a *Area) NumsImage(title string) image.Image {
	return WriteOnGrid(lo.Keys(*a), title, func(region Region) string {
		return fmt.Sprintf("%v", (*a)[region])
	})
}

func (a *Area) Image(imageType int, title string) image.Image {
	switch imageType {
	case PointsImage:
		return a.PointsImage(title)
	case NumsImage:
		return a.NumsImage(title)
	default:
		panic("unknown type")
	}
}

func (a Area) MarshalJSON() ([]byte, error) {
	ln, i := len(a), 0

	buf := bytes.Buffer{}
	buf.WriteString("{")

	for reg, val := range a {
		key, err := reg.MarshalJSON()
		if err != nil {
			return nil, err
		}

		buf.WriteString(`"`)
		buf.Write(key)
		buf.WriteString(`"`)
		buf.WriteString(":")
		buf.WriteString(strconv.Itoa(val))

		if i != ln-1 {
			buf.Write([]byte(","))
		}

		i++
	}

	buf.WriteString("}")
	return buf.Bytes(), nil
}

func (a *Area) UnmarshalJSON(data []byte) error {
	if (*a) == nil {
		*a = map[Region]int{}
	}

	tmpMap := map[string]json.RawMessage{}

	err := json.Unmarshal(data, &tmpMap)
	if err != nil {
		return err
	}

	for key, val := range tmpMap {
		r := Region{}

		err = r.UnmarshalJSON([]byte(key))
		if err != nil {
			return err
		}

		v, err := strconv.Atoi(string(val))
		if err != nil {
			return err
		}

		(*a)[r] = v
	}

	return nil
}
