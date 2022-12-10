package consensus

type Consensus struct {
	SnowParams
}

type SnowParams struct {
	K int // K sample K of each round of query. K < number_of_peers
	A int // A number of threshold that can be considered a majority. A < K
	B int // B number of rounds of successive agreement of sample K
}

func NewConsensus(params SnowParams) *Consensus {
	return &Consensus{
		SnowParams: params,
	}
}
