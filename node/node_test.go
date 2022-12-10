package node

import (
	"testing"
)

func TestNewNode(t *testing.T) {

	// new node
	node := NewNode("127.0.0.1:9447")
	node.StartServer()
	defer node.StopServer()
}

func TestDiscoverPeers(t *testing.T) {

	// new node
	node1 := NewNode("127.0.0.1:9447")
	node1.StartServer()

	node2 := NewNode("127.0.0.1:9448")
	node2.StartServer()

	// connect to node2
	node1.PeerManager.StartDiscoverPeers(node2.Address)
	node1.Waiter.Wait()

}
