package node

import "testing"

func TestNewNode(t *testing.T) {

	// new node
	node := NewNode("127.0.0.1:9447")
	node.StartServer()
	defer node.StopServer()
}
