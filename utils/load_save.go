package utils

import (
	"encoding/json"
	"gonum.org/v1/gonum/mat"
	"io"
	"os"
)

type Marshal func(v any) ([]byte, error)
type Unmarshal func(data []byte, v any) error

func MustSave[T any](path string, d T, m Marshal) {
	data, err := m(d)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	_, err = file.Write(data)
	if err != nil {
		panic(err)
	}
}

func MustLoad[T any](path string, d *T, u Unmarshal) {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	err = u(data, d)
	if err != nil {
		panic(err)
	}
}

func DenseMarshal(dense *mat.Dense) Marshal {
	return func(v any) ([]byte, error) {
		return dense.MarshalBinary()
	}
}

func ManyDenseMarshal(dense []*mat.Dense) Marshal {
	return func(v any) ([]byte, error) {
		buf := [][]byte{}

		for _, d := range dense {
			tmp, err := d.MarshalBinary()
			if err != nil {
				return nil, err
			}

			buf = append(buf, tmp)
		}

		return json.Marshal(buf)
	}
}

func DenseUnmarshal(dense *mat.Dense) Unmarshal {
	return func(data []byte, v any) error {
		return dense.UnmarshalBinary(data)
	}
}

func ManyDenseUnmarshal(dense *[]*mat.Dense) Unmarshal {
	return func(data []byte, v any) error {
		buf := [][]byte{}

		if err := json.Unmarshal(data, &buf); err != nil {
			return err
		}

		for _, b := range buf {
			nd := &mat.Dense{}
			if err := nd.UnmarshalBinary(b); err != nil {
				return err
			}

			*dense = append(*dense, nd)
		}

		return nil
	}
}
