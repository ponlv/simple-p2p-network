package node

import (
	"google.golang.org/grpc"
	"simple-p2p/p2p"
	"sync"
)

// Node represents a P2P node in the network.
type Node struct {
	// Config fields may not be modified while the node is running.
	Address string // Network address of the node

	Server *grpc.Server // gRPC server instance

	Waiter *sync.WaitGroup // WaitGroup for graceful shutdown

	PeerManager p2p.PeerManager // Peer manager instance
}

// NewNode creates a new node instance.
func NewNode(address string) *Node {
	return &Node{
		Address:     address,
		Server:      grpc.NewServer(),
		Waiter:      &sync.WaitGroup{},
		PeerManager: p2p.NewPeerManager(address),
	}
}
