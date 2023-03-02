package utils

import (
	"fmt"
	"math"
)

const FloatPrecision = 4

func SetPrecision(n float64, precision int) float64 {
	tmp := math.Pow(10, float64(precision))

	return math.Round(n*tmp) / tmp
}

func PrecisionString(n float64) string {
	return fmt.Sprintf("%v", SetPrecision(n, FloatPrecision))
}

func WrapInBox(str string, cellWidth int) string {
	dif := cellWidth - len(str)
	sp1 := Spaces(dif / 2)
	sp2 := Spaces(dif/2 + dif%2)

	return "|" + sp1 + str + sp2 + "|"
}

func Spaces(count int) string {
	spaces := ""

	for j := 0; j < count; j++ {
		spaces += " "
	}

	return spaces
}
