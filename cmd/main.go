package main

import (
	"diploma/markov"
	"diploma/poisson"
	"fmt"
	"gonum.org/v1/gonum/mat"
	"math"
)

const (
	xStartG           = 0
	xEndG             = 10
	yStartG           = 0
	yEndG             = 10
	XAxisPartQuantity = 2
	YAxisPartQuantity = 2

	xCenter = (xEndG - xStartG) / 2
	yCenter = (yEndG - yStartG) / 2

	N = 50
)

func regular(x, y float64) float64 {
	return math.Max(-math.Pow(x-y, 2)/2+30, 0)
}

func irregular(x, y float64) float64 {
	return math.Max(math.Pow(y-yCenter, 3)+10, 4)
}

func cand1(x, y float64) float64 {
	return y * y / 10
}

func cand2(x, y float64) float64 {
	return math.Pow(math.E, 4+(y-x)*(x-y)/24) / 15
}

func main() {
	p1 := poisson.NewProcessWithIntensityFunc(xStartG, xEndG, yStartG, yEndG, XAxisPartQuantity, YAxisPartQuantity, cand1)
	p2 := poisson.NewProcessWithIntensityFunc(xStartG, xEndG, yStartG, yEndG, XAxisPartQuantity, YAxisPartQuantity, cand2)

	mu := mat.NewDense(1, 2, []float64{0.75, 0.25})
	a := mat.NewDense(2, 2, []float64{
		0.8, 0.2,
		0.5, 0.5,
	})

	proc := markov.NewGenerator(mu, a, []*poisson.ProcessWithIntensityFunc{p1, p2})

	frames := proc.Generate(N)

	muL := mat.NewDense(1, 2, []float64{0.70, 0.30})
	aL := mat.NewDense(2, 2, []float64{
		0.65, 0.35,
		0.45, 0.55,
	})
	process := []map[poisson.Region]float64{p2.IntensityMap(), p1.IntensityMap()}
	fmt.Println("1 mmm", p1.IntensityMap())
	fmt.Println("2 mmm", p2.IntensityMap())
	l := markov.NewLearner(muL, aL, process, frames)

	for i := 0; i < 30; i++ {
		//fmt.Println("---------------- new step")
		l.Step()
	}

	fmt.Println("mu:", l.Mu)
	fmt.Println("a", l.A)
	fmt.Println("process 1", l.Processes[0])
	fmt.Println("process 2", l.Processes[1])
	for i, frame := range frames {
		drawFrame(frame, i)
	}
	//gen := p2.Generate()
	//fmt.Println(gen.TotalPoint())
	//gen.Draw("card.png")
	//
	//im := p1.IntensityMap()
	//
	//pl := poisson.DrawGrid(lo.Keys(im), "real intensity")
	//if err := pl.Save(4*vg.Inch, 4*vg.Inch, "tmp.png"); err != nil {
	//	panic(err)
	//}
	//
	//file, err := os.OpenFile("tmp.png", os.O_RDWR, 0666)
	//if err != nil {
	//	panic(err)
	//}
	//
	//imgPure, err := png.Decode(file)
	//if err != nil {
	//	panic(err)
	//}
	//
	//img := utils.ImageToRGBA(imgPure)
	//
	//min, max := img.Rect.Min, img.Rect.Max
	//
	//sizeX, sizeY := max.X-min.X, max.Y-min.Y
	//relSizeX, relSizeY := float64(sizeX)/10, float64(sizeY)/10
	//relSizeX, relSizeY = 32, 31
	//addX, addY := 43.0, 30.0
	//
	//for region, val := range im {
	//	x, y := region.Center()
	//	y = 10 - y
	//
	//	utils.AddLabel(img, int(x*relSizeX+addX), int(y*relSizeY+addY), formatNum(val))
	//}
	//
	//f, _ := os.OpenFile("intensity.png", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	//fmt.Println(png.Encode(f, img))
}

func drawFrame(area *poisson.Area, frameIndex int) {
	area.Draw(fmt.Sprintf("./tmp/card-%v.png", frameIndex))
}

func formatNum(num float64) string {
	if num > 100 {
		return fmt.Sprintf("%3.0f", num)
	}

	if num > 10 {
		return fmt.Sprintf("%3.1f", num)
	}

	return fmt.Sprintf("%3.2f", num)
}
