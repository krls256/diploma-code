package markov

import (
	"diploma/poisson"
)

type ResultChain struct {
	Frames     []*poisson.Area `json:"ObservationSequences"`
	StateChain []int           `json:"state_chain"`
}
