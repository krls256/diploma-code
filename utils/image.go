package utils

import (
	"bytes"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"image"
	"image/color"
	"image/draw"
	"image/png"
)

func AddLabel(img *image.RGBA, x, y int, label string) {
	col := color.Black

	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
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
	w, err := p.WriterTo(4*vg.Inch, 4*vg.Inch, "png")
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
