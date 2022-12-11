package consensus

import "simple-p2p/node"

type Consensus struct {
	SnowParams
	Node node.Node

	preference int  // preference of the node
	confident  int  // confidence of the node
	accepted   bool // accepted value of the node
}

type SnowParams struct {
	K       int // K sample K of each round of query. K < number_of_peers
	A       int // A is quorum size. A < K
	B       int // B is decision threshold
	MaxStep int // MaxStep is the maximum number of rounds of query
}

// NewConsensus creates a new consensus instance.
func NewConsensus(params SnowParams) *Consensus {
	return &Consensus{
		SnowParams: params,
	}
}

// Preference returns the preference of the node.
func (c *Consensus) Preference() int {
	return c.preference
}

// Start starts the consensus process.
func (c *Consensus) Start() {

	c.confident = 1
	c.accepted = false

	for ; c.accepted == false; {
		
	}
	// TODO: get k peers from the peer manager
	kPeers := c.Node.GetSamplePeers(c.K)

	// TODO: send query to each peer
	responses := make([]int, c.K)
	for _, peer := range kPeers {
		responses = append(responses, c.Node.GetPreference(peer))
	}

	// TODO: get most frequent value from responses
	count, value := c.GetMostFrequentValue(responses)

	// TODO: check if frequency is greater than A
	if count > c.A {
		oldPreference := c.preference

		c.preference = value

		if oldPreference != c.preference {
			c.confident = 1
		} else {
			c.confident++

			if c.confident > c.B {
				c.accepted = true
			}
		}
	}
}
