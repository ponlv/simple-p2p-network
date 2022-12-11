package node

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var host = "127.0.0.1"

func TestDiscoverPeers(t *testing.T) {
	
	node1 := createNode(9447)
	node1.StartServer()

	node2 := createNode(9448)
	node2.StartServer()

	node3 := createNode(9449)
	node3.StartServer()

	node4 := createNode(9450)
	node4.StartServer()

	node5 := createNode(9451)
	node5.StartServer()

	node6 := createNode(9452)
	node6.StartServer()

	node1.PeerManager.StartDiscoverPeers(node2.Address, node3.Address, node4.Address, node5.Address, node6.Address)
	node2.PeerManager.StartDiscoverPeers(node1.Address)
	node3.PeerManager.StartDiscoverPeers(node1.Address)
	node4.PeerManager.StartDiscoverPeers(node1.Address)
	node5.PeerManager.StartDiscoverPeers(node1.Address)
	node6.PeerManager.StartDiscoverPeers(node1.Address)

	// wait for each node to discover all others nodes
	time.Sleep(2 * time.Second)

	assert.Equal(t, 5, len(node1.PeerManager.GetPeers()))
	assert.Equal(t, 5, len(node2.PeerManager.GetPeers()))
	assert.Equal(t, 5, len(node3.PeerManager.GetPeers()))
	assert.Equal(t, 5, len(node4.PeerManager.GetPeers()))
	assert.Equal(t, 5, len(node5.PeerManager.GetPeers()))
	assert.Equal(t, 5, len(node6.PeerManager.GetPeers()))
}

func createNode(port int64) *Node {
	return NewNode(fmt.Sprintf("%v:%d", host, port))
}
