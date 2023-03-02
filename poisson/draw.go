package poisson

import (
	"diploma/utils"
	"fmt"
	"github.com/andybons/gogif"
	"github.com/samber/lo"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"os"
)

const (
	PointsImage = iota
	NumsImage
)

func (a *Area) PointsImage(titleNum int) image.Image {
	scatterData := a.Points()

	p := DrawGrid(lo.Keys(*a), fmt.Sprintf("Points Process %v", titleNum))

	s, err := plotter.NewScatter(scatterData)
	if err != nil {
		panic(err)
	}
	s.Radius = vg.Points(1)
	s.GlyphStyle.Color = color.RGBA{A: 255}

	p.Legend.Add("scatter", s)

	p.Add(s)

	return utils.PlotToPNG(p)
}

func (a *Area) NumsImage(titleNum int) image.Image {
	return WriteOnGrid(lo.Keys(*a), fmt.Sprintf("Points Process %v", titleNum), func(region Region) string {
		return fmt.Sprintf("%v", (*a)[region])
	})
}

func (a *Area) Image(imageType, titleNum int) image.Image {
	switch imageType {
	case PointsImage:
		return a.PointsImage(titleNum)
	case NumsImage:
		return a.NumsImage(titleNum)
	default:
		panic("unknown type")
	}
}

func DrawIntensityMap(im map[Region]float64, titleNum int) image.Image {
	return WriteOnGrid(lo.Keys(im), fmt.Sprintf("Points Process %v", titleNum), func(region Region) string {
		return fmt.Sprintf("%.3f", im[region])
	})
}

func DrawFrames(areas []*Area, filenameBase string, imageType int) {
	for i, area := range areas {
		filename := fmt.Sprintf("%v-%v.png", filenameBase, i+1)
		img := area.Image(imageType, i+1)

		f, err := os.OpenFile(filename, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}

		if err := png.Encode(f, img); err != nil {
			panic(err)
		}
	}
}

func DrawGif(areas []*Area, stageChain []int, intensityMaps []map[Region]float64, filenameBase string) {
	f, err := os.OpenFile(fmt.Sprintf("%v.gif", filenameBase), os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	anim := &gif.GIF{}

	for i, area := range areas {
		currencuIntensityMap := intensityMaps[stageChain[i]]

		img1 := area.Image(NumsImage, i+1)
		img2 := area.Image(PointsImage, i+1)
		img3 := DrawIntensityMap(currencuIntensityMap, i+1)

		img := utils.HorizontalJoinImage(img1, img2, img3)

		bounds := img.Bounds()
		palettedImage := image.NewPaletted(bounds, nil)

		quantizer := gogif.MedianCutQuantizer{NumColor: 64}
		quantizer.Quantize(palettedImage, bounds, img, image.ZP)

		// Add new frame to animated GIF
		anim.Image = append(anim.Image, palettedImage)
		anim.Delay = append(anim.Delay, 100)
	}

	if err := gif.EncodeAll(f, anim); err != nil {
		panic(err)
	}
}

func WriteOnGrid(regions []Region, imageTitle string, textFunc func(Region) string) image.Image {
	p := DrawGrid(regions, imageTitle)

	imgPure := utils.PlotToPNG(p)
	img := utils.ImageToRGBA(imgPure)

	min, max := img.Rect.Min, img.Rect.Max

	sizeX, sizeY := max.X-min.X, max.Y-min.Y
	relSizeX, relSizeY := float64(sizeX)/10, float64(sizeY)/10
	relSizeX, relSizeY = 32, 31
	addX, addY := 43.0, 30.0

	for _, region := range regions {
		x, y := region.Center()
		y = 10 - y

		utils.AddLabel(img, int(x*relSizeX+addX), int(y*relSizeY+addY), textFunc(region))
	}

	return img
}