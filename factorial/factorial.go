package factorial

import "errors"

const (
	// OverflowError is the error message provided when the result of a factorial computation overflows uint64
	OverflowError = "overflow error"
	// Negative is the error provided when the Factorial given input is negative
	Negative = "input is below 0"
)

func Factorial(n int) (uint64, error) {
	if n < 0 {
		return 0, errors.New(Negative)
	}

	if n >= 21 {
		return 0, errors.New(OverflowError)
	}

	if n < 2 {
		return 1, nil
	}

	un := uint64(n)
	m := uint64(1)

	for ; un > 1; un-- {
		m = m * un
	}

	return m, nil
}
