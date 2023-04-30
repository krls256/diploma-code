package poisson

import (
	"bytes"
	"diploma/utils"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"golang.org/x/exp/rand"
	"sort"
	"strconv"
)

type IntensityMap map[Region]float64

func (im IntensityMap) Add(im1, im2 IntensityMap) {
	for reg := range im2 {
		im[reg] = im1[reg] + im2[reg]
	}
}

func (im IntensityMap) Typeless() map[[4]float64]float64 {
	m := map[[4]float64]float64{}

	for k, v := range im {
		m[k] = v
	}

	return m
}

func (im IntensityMap) String() string {
	buf := bytes.Buffer{}

	regions := lo.Keys(im)
	sort.Slice(regions, func(i, j int) bool {
		return RegionRowOrderBiggerThan(regions[i], regions[j])
	})

	for i := 0; i < len(regions); i++ {
		if i > 0 && regions[i-1][3] != regions[i][3] {
			buf.WriteString("\n")
		}
		str := utils.PrecisionString(im[regions[i]])

		buf.WriteString(utils.WrapInBox(str, StringCellWidth))
	}

	return buf.String()
}

func (im IntensityMap) Copy() IntensityMap {
	newIm := IntensityMap{}

	for k, v := range im {
		newIm[k] = v
	}

	return newIm
}

func (im IntensityMap) Generate() *Area {
	area := Area{}

	intensity := 0.0

	for _, in := range im {
		intensity += in
	}

	process := poisson(intensity)
	points := int(process.Rand())

	for i := 0; i < points; i++ {
		r := rand.Float64() * intensity

		for reg, in := range im {
			if _, ok := area[reg]; !ok {
				area[reg] = 0
			}

			if r < in {
				area[reg]++

				break
			}

			r -= in
		}
	}

	return &area
}

func (a IntensityMap) MarshalJSON() ([]byte, error) {
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
		buf.WriteString(fmt.Sprintf("%f", val))

		if i != ln-1 {
			buf.Write([]byte(","))
		}

		i++
	}

	buf.WriteString("}")
	return buf.Bytes(), nil
}

func (a *IntensityMap) UnmarshalJSON(data []byte) error {
	if (*a) == nil {
		*a = map[Region]float64{}
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

		v, err := strconv.ParseFloat(string(val), 64)
		if err != nil {
			return err
		}

		(*a)[r] = v
	}

	return nil
}

func StandardDeviation(a, b IntensityMap) float64 {
	sd, tmp := 0.0, 0.0

	for reg := range a {
		tmp = a[reg] - b[reg]

		sd += tmp * tmp
	}

	return sd
}
