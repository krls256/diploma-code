package utils

import (
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

func DenseUnmarshal(dense *mat.Dense) Unmarshal {
	return func(data []byte, v any) error {
		return dense.UnmarshalBinary(data)
	}
}
