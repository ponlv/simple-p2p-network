package node

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"simple-p2p/p2p"
	"simple-p2p/p2p/message"
	"simple-p2p/proto/proto"
	"sync"
	"time"
)

// Node represents a P2P node in the network.
type Node struct {
	// Config fields may not be modified while the node is running.
	Address string // Network address of the node

	Server *grpc.Server // gRPC server instance

	Waiter *sync.WaitGroup // WaitGroup for graceful shutdown

	PeerManager p2p.Peer // Peer manager instance

	MessageManager message.MessageManager // Message manager instance

}

// NewNode creates a new node instance.
func NewNode(address string) *Node {
	return &Node{
		Address:        address,
		Server:         grpc.NewServer(),
		Waiter:         &sync.WaitGroup{},
		PeerManager:    p2p.NewPeerManager(address),
		MessageManager: message.NewMessageManager(),
	}
}

// StartServer starts server to provide services. This must be called after
// registering any other external service.
func (n *Node) StartServer() {
	conn, _ := net.DialTimeout("tcp", n.Address, 5*time.Second)
	if conn != nil {
		conn.Close()
	}

	lis, err := net.Listen("tcp", n.Address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// register internal service
	proto.RegisterPeerServiceServer(n.Server, n.PeerManager)
	proto.RegisterMessageServiceServer(n.Server, n.MessageManager)

	log.Printf("server is listening at: %v", n.Address)
	n.Waiter.Add(1)

	go func() {
		err := n.Server.Serve(lis)
		if err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}

// StopServer stops the server.
func (n *Node) StopServer() {
	if n.Server != nil {
		n.Server.Stop()
		log.Printf("server stopped: %v", n.Address)
		n.Waiter.Done()
	}
}
