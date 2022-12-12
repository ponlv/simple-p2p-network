package main

import (
	"flag"
	"fmt"
	"simple-p2p/consensus"
	"simple-p2p/node"
	"simple-p2p/proto/proto"
)

func main() {

	// add flag
	neighbors := flag.String("neighbors", "", "bootstrap address to join the p2p network")
	host := flag.String("host", "127.0.0.1", "host address")
	port := flag.Int64("port", 9447, "port to listen")
	K := flag.Int("K", 3, "sample K of each round of query. K < number_of_peers")
	Alpha := flag.Int("A", 2, "is quorum size. A < K")
	Beta := flag.Int("B", 10, "is decision threshold")
	MaxStep := flag.Int("max-step", 100, "is the maximum number of rounds of query")
	flag.Parse()

	// start node
	newNode := node.NewNode(fmt.Sprintf("%v:%d", *host, *port))

	// start peer discovery
	if *neighbors != "" {
		newNode.PeerManager.StartDiscoverPeers(*neighbors)
	}

	// start consensus
	snow := consensus.NewConsensus(consensus.SnowParams{
		K:       *K,
		A:       *Alpha,
		B:       *Beta,
		MaxStep: *MaxStep,
	})
	snow.AddNode(newNode)
	proto.RegisterConsensusServiceServer(newNode.Server, snow)

	//  start server
	newNode.StartServer()

	newNode.Waiter.Wait()

	fmt.Printf("Node %v is started \n", newNode.Address)
}
