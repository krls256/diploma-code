package main

import (
	"fmt"
	"gitlab.eemi.tech/golang/factgorial/factgorial"
	"math"
)

var (
	y = 4.0
	g = 5.625
)

func main() {
	f, _ := factgorial.Factorial(int(y))
	fmt.Println(math.Pow(g, y) / float64(f) * math.Pow(math.E, -g))
}
