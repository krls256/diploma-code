package constants

import "fmt"

const (
	CacheDir                   = "./cache"
	ChainFile                  = "chain.json"
	BaseDistributionFile       = "base-distribution.bin"
	HiddenDistributionFile     = "hidden-distribution.bin"
	ObservableDistributionFile = "observable-distribution.json"
	IntensityMapFormat         = "intensity-%v-%v.json"
)

func CachePath(filename string) string {
	return fmt.Sprintf("%s/%s", CacheDir, filename)
}
