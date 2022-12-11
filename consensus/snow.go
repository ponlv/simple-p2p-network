package consensus

import (
	"context"
	"fmt"
	"simple-p2p/node"
	"simple-p2p/proto/proto"
	"simple-p2p/utils"
	"sync"
)

type Consensus interface {

	// Preference returns the preference of the node.
	Preference() int

	// Sync starts the consensus process.
	Sync()

	// GetPreference is internal call to perform a single step of the consensus
	GetPreference(context.Context, *proto.Empty) (*proto.GetPreferenceResponse, error)

	// AddNode adds a node to the consensus.
	AddNode(*node.Node)

	// GetNode returns the node of the consensus.
	GetNode() *node.Node

	// UpdatePreference updates the preference of the node.
	UpdatePreference(int)
}

var _ Consensus = (*consensus)(nil)

type consensus struct {
	SnowParams
	Node *node.Node

	preference int          // preference of the node
	confident  int          // confidence of the node
	accepted   bool         // accepted value of the node
	isRunning  bool         // consensus is running
	mux        sync.RWMutex // mutual exclusion lock for peers
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

// Sync starts the consensus process.
func (c *consensus) Sync() {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.confident = 1
	c.accepted = false

	i := 0
	for ; c.accepted == false; i++ {
		fmt.Printf("Node %v: Round %d: preference = %d, confident = %d, accepted = %t \n", c.Node.Address, i, c.preference, c.confident, c.accepted)

		c.step()

		if i > c.MaxStep {
			fmt.Printf("Node %v: Consensus failed \n", c.Node.Address)
			break
		}
	}

	fmt.Printf("Node %v: Consensus succeeded after %v rounds \n", c.Node.Address, i)
}

// step performs a single step of the consensus.
func (c *consensus) step() {
	// get K peers from the peer manager
	kPeers := c.Node.PeerManager.GetSamplePeers(c.K)

	// send query to each peer
	responses := make([]int, c.K)
	for i, peer := range kPeers {
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

		responses[i] = int(response.Preference)
	}

	// get most frequent value from responses
	value, count := utils.GetMostFrequentValue(responses)

	// check if frequency is greater than A
	if count >= c.A {
		oldPreference := c.preference

		c.preference = value

		// check if preference is changed, the confidence is reset to 1
		// otherwise, the confidence is increased by 1
		if oldPreference != c.preference {
			c.confident = 1
		} else {
			c.confident++

			// check if confidence is greater than B, the value is accepted
			if c.confident >= c.B {
				c.accepted = true
			}
		}
	} else {
		c.confident = 0
	}

}

// GetPreference returns the preference of the node.
func (c *consensus) GetPreference(context.Context, *proto.Empty) (*proto.GetPreferenceResponse, error) {
	return &proto.GetPreferenceResponse{
		Preference: int64(c.preference),
	}, nil
}

// AddNode adds a node to the consensus.
func (c *consensus) AddNode(n *node.Node) {
	c.Node = n
}

// GetNode returns the node of the consensus.
func (c *consensus) GetNode() *node.Node {
	return c.Node
}

// UpdatePreference updates the preference of the node.
func (c *consensus) UpdatePreference(p int) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.preference = p
}
