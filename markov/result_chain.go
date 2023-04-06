package markov

import (
	"diploma/poisson"
)

type ResultChain struct {
	Frames     []*poisson.Area `json:"Frames"`
	StateChain []int           `json:"state_chain"`
}
