package consensus

import (
	"fmt"
	"simple-p2p/node"
	"simple-p2p/proto/proto"
	"testing"
	"time"
)

var host = "127.0.0.1"

func TestSnow(t *testing.T) {
	var choices = []int{1, 2, 3}

	t.Run("TestSnow", func(t *testing.T) {
		numNode := 5
		listConsensus := make([]Consensus, numNode)
		for i := 0; i < numNode; i++ {
			newNode := createNode(9447 + int64(i))

			consensus := NewConsensus(SnowParams{
				K:       3,
				A:       2,
				B:       10,
				MaxStep: 100,
			})
			consensus.AddNode(newNode)

			// pick a random choice and update the preference
			consensus.UpdatePreference(choices[i%3])

			// create a new consensus instance
			proto.RegisterConsensusServiceServer(newNode.Server, consensus)
			listConsensus[i] = consensus

			newNode.StartServer()
			fmt.Printf("Node %v: Started\n", newNode.Address)
		}

		// connect all nodes
		for i := 1; i < numNode; i++ {
			listConsensus[i].GetNode().PeerManager.StartDiscoverPeers(listConsensus[i-1].GetNode().Address)
		}

		// wait for each node to discover all others nodes
		time.Sleep(15 * time.Second)

		for i := 0; i < numNode; i++ {
			fmt.Printf("Node %v: %v\n", listConsensus[i].GetNode().Address, len(listConsensus[i].GetNode().PeerManager.GetPeers()))
		}

		// Start the consensus
		for i := 0; i < numNode; i++ {
			go listConsensus[i].Sync()
		}

		// wait for the consensus to finish
		time.Sleep(10 * time.Second)

		// check if the consensus has the same preference
		for i := 1; i < numNode; i++ {
			if listConsensus[i].Preference() != listConsensus[i-1].Preference() {
				t.Errorf("Consensus failed")
			}
		}
	})
}

func createNode(port int64) *node.Node {
	return node.NewNode(fmt.Sprintf("%v:%d", host, port))
}
