package utils

import (
	"bytes"
	"diploma/constants"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image"
	"image/color"
	"image/draw"
	"image/png"
)

func AddLabel(img *image.RGBA, x, y int, labels ...string) {
	col := color.Black

	baseY := y - (constants.LineHeightPixel*(len(labels)-1))/2

	for _, label := range labels {
		point := fixed.Point26_6{X: fixed.I(x - constants.SymbolWidthPixel*len(label)/4), Y: fixed.I(baseY)}

		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(col),
			Face: basicfont.Face7x13,
			Dot:  point,
		}

		d.DrawString(label)

		baseY += constants.LineHeightPixel
	}
}

func FillColor(img *image.RGBA, col color.Color) {
	x, y := img.Rect.Dx(), img.Rect.Dy()

	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			img.Set(i, j, col)
		}
	}
}

func ImageToRGBA(src image.Image) *image.RGBA {
	if dst, ok := src.(*image.RGBA); ok {
		return dst
	}

	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
	return dst
}

func HorizontalJoinImage(im image.Image, imgLeft ...image.Image) image.Image {
	imgs := append([]image.Image{im}, imgLeft...)

	point := image.Point{X: imgs[0].Bounds().Dx()}
	rct := image.Rectangle{Min: point, Max: point}
	rcts := []image.Rectangle{{
		Min: image.Point{0, 0},
		Max: image.Point{X: imgs[0].Bounds().Dx(), Y: imgs[0].Bounds().Dy()}}}

	for i, im := range imgs {
		if i != 0 {
			prevRcts := rcts[len(rcts)-1]
			min := prevRcts.Min
			min.X = min.X + im.Bounds().Dx()
			max := prevRcts.Max
			max.X = max.X + im.Bounds().Dx()

			rcts = append(rcts, image.Rectangle{Min: min, Max: max})
			rct = image.Rectangle{Min: rct.Min, Max: max}
		}
	}

	r := image.Rectangle{Min: image.Point{0, 0}, Max: rct.Max}
	img := image.NewRGBA(r)

	for i, im := range imgs {
		draw.Draw(img, rcts[i], im, image.Point{0, 0}, draw.Src)
	}

	return img
}

func PlotToPNG(p *plot.Plot) image.Image {
	w, err := p.WriterTo(constants.ImageInchWidth, constants.ImageInchHeight, "png")
	if err != nil {
		panic(err)
	}

	buf := &bytes.Buffer{}

	if _, err = w.WriteTo(buf); err != nil {
		panic(err)
	}

	imgPure, err := png.Decode(buf)
	if err != nil {
		panic(err)
	}

	return imgPure
}

func DrawValuesGraphics(values []float64, title, filename string) {
	p := plot.New()

	p.Title.Text = title
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(plotter.NewGrid())

	lineData := plotter.XYs{}
	for i, pr := range values {
		lineData = append(lineData, plotter.XY{X: float64(i), Y: pr})
	}

	l, err := plotter.NewLine(lineData)
	if err != nil {
		panic(err)
	}
	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	l.LineStyle.Color = color.RGBA{B: 255, A: 255}

	p.Add(l)

	if err = p.Save(constants.ImageInchWidth, constants.ImageInchHeight, filename+".png"); err != nil {
		panic(err)
	}
}

func DrawVariablesProgress(vars [][]float64, legengs []string, title, xName, yName, filename string) {
	p := plot.New()

	p.Title.Text = title
	p.X.Label.Text = xName
	p.Y.Label.Text = yName

	c := colorer()

	for i := 0; i < len(vars); i++ {
		lineData := plotter.XYs{}
		for i, pr := range vars[i] {
			lineData = append(lineData, plotter.XY{X: float64(i), Y: pr})
		}

		l, err := plotter.NewLine(lineData)
		if err != nil {
			panic(err)
		}
		l.LineStyle.Width = vg.Points(1)
		l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
		l.LineStyle.Color = c()

		p.Add(l)
		p.Legend.Add(legengs[i], l)
	}

	if err := p.Save(constants.ImageInchWidth, constants.ImageInchHeight, filename+".png"); err != nil {
		panic(err)
	}
}

func colorer() func() color.RGBA {
	init := color.RGBA{R: 0, G: 85, B: 170, A: 255}
	return func() color.RGBA {
		init.R += 86
		init.G += 86
		init.B += 86
		return init
	}
}
