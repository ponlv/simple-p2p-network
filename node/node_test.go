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

func TestSendMessage(t *testing.T) {
	//// new node
	//node1 := NewNode("127.0.0.1:9447")
	//node1.StartServer()
	//
	//node2 := NewNode("127.0.0.1:9448")
	//node2.StartServer()
	//
	//node2Conn, err := node1.PeerManager.GetConnection(node2.Address)
	//if err != nil {
	//	return
	//}
	//
	//err = node1.MessageManager.SendMessage(node2Conn, &proto.MessageRequest{
	//	Type:  proto.MessageType_VALUE,
	//	Value: []byte("hello"),
	//})
	//if err != nil {
	//	return
	//}

}
