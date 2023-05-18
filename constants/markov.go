package constants

const (
	FieldXStart = 0
	FieldXEnd   = 10
	FieldYStart = 0
	FieldYEnd   = 10

	MaxPoissonIntensity = 9
	IntensityStep       = 0.01

	FieldXCenter = (FieldXEnd - FieldXStart) / 2
	FieldYCenter = (FieldYEnd - FieldYStart) / 2

	ObservationSequences = 10
	LearnersCount        = 1

	ItersPerEpoch = 5
	SampleSize    = 250
)

var (
	XAxisSplit = 2
	YAxisSplit = 2
)

func AddSplit() {
	XAxisSplit++
	YAxisSplit++
}
