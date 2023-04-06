package poisson

import (
	"fmt"
	"strconv"
	"strings"
)

type Region [4]float64 // XStart, XEnd, YStart, YEnd

func (r *Region) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("(%f,%f,%f,%f)", r[0], r[1], r[2], r[3])), nil
}

func (r *Region) UnmarshalJSON(data []byte) error {
	data = data[1 : len(data)-1]

	str := string(data)
	nums := strings.Split(str, ",")

	for i, num := range nums {
		f, err := strconv.ParseFloat(num, 64)
		if err != nil {
			return err
		}

		r[i] = f
	}

	return nil
}

func RegionRowOrderBiggerThan(r1, r2 Region) bool {
	if r1[2] == r2[2] {
		return r1[0] > r2[0]
	}

	return r1[2] > r2[2]
}

func (r Region) String() string {
	return fmt.Sprintf("%v-%v-%v-%v", r[0], r[1], r[2], r[3])
}

func (r Region) Center() (xCenter, yCenter float64) {
	xStart, xEnd, yStart, yEnd := r[0], r[1], r[2], r[3]
	return xStart + (xEnd-xStart)/2, yStart + (yEnd-yStart)/2
}
