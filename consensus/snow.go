package consensus

import (
	"context"
	"simple-p2p/node"
	"simple-p2p/proto/proto"
	"simple-p2p/utils"
)

type Consensus interface {
	Preference() int
	Sync()
	GetPreference(context.Context, *proto.Empty) (*proto.GetPreferenceResponse, error)
}

var _ Consensus = (*consensus)(nil)

type consensus struct {
	SnowParams
	Node node.Node

	preference int  // preference of the node
	confident  int  // confidence of the node
	accepted   bool // accepted value of the node
	isRunning  bool // consensus is running
}

type SnowParams struct {
	K       int // K sample K of each round of query. K < number_of_peers
	A       int // A is quorum size. A < K
	B       int // B is decision threshold
	MaxStep int // MaxStep is the maximum number of rounds of query
}

// NewConsensus creates a new consensus instance.
func NewConsensus(params SnowParams) Consensus {
	return &consensus{
		SnowParams: params,
	}
}

// Preference returns the preference of the node.
func (c *consensus) Preference() int {
	return c.preference
}

// Start starts the consensus process.
func (c *consensus) Sync() {

	if c.isRunning {
		return
	}

	c.confident = 1
	c.accepted = false
	c.isRunning = true

	for c.accepted == false {

		c.step()
	}

	c.isRunning = false
}

func (c *consensus) step() {
	// get K peers from the peer manager
	kPeers := c.Node.PeerManager.GetSamplePeers(c.K)

	// send query to each peer
	responses := make([]int, c.K)
	for _, peer := range kPeers {
		// get connection of the peer
		conn, err := c.Node.PeerManager.GetConnection(peer)
		if err != nil {
			continue
		}

		// create a proto client
		client := proto.NewConsensusServiceClient(conn)

		// send query
		response, err := client.GetPreference(context.Background(), &proto.Empty{})
		if err != nil {
			continue
		}

		responses = append(responses, int(response.Preference))
	}

	// get most frequent value from responses
	count, value := utils.GetMostFrequentValue(responses)

	// check if frequency is greater than A
	if count > c.A {
		oldPreference := c.preference

		c.preference = value

		// check if preference is changed, the confidence is reset to 1
		// otherwise, the confidence is increased by 1
		if oldPreference != c.preference {
			c.confident = 1
		} else {
			c.confident++

			// check if confidence is greater than B, the value is accepted
			if c.confident > c.B {
				c.accepted = true
			}
		}
	} else {
		c.confident = 0
	}
}

func (c *consensus) GetPreference(context.Context, *proto.Empty) (*proto.GetPreferenceResponse, error) {
	return &proto.GetPreferenceResponse{
		Preference: int64(c.preference),
	}, nil
}
