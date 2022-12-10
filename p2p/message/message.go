package message

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"google.golang.org/grpc"
	"simple-p2p/proto/proto"
	"time"
)

type MessageManager interface {

	// SendMessage sends a message to a peer.
	SendMessage(conn *grpc.ClientConn, message *proto.MessageRequest) error

	// ReceiveMessage receives a message from a peer.
	ReceiveMessage(context.Context, *proto.MessageRequest) (*proto.MessageResponse, error)
}

var _ MessageManager = (*messageManager)(nil)

// messageLog is a log item for a message. Only one of sender and receiver
// need to be assigned.
type messageLog struct {
	hash        string
	messageType int
	sender      string
	receiver    string
	time        time.Time
}

// MessageManager is the service to receive and process messages.
type messageManager struct {
	MessageLogs []messageLog // logs for sent/received messages
}

// NewMessageManager creates a new message manager instance.
func NewMessageManager() MessageManager {
	return &messageManager{
		MessageLogs: make([]messageLog, 0),
	}
}

// SendMessage sends a message to a peer with given grpc connection.
func (m *messageManager) SendMessage(conn *grpc.ClientConn, message *proto.MessageRequest) error {
	// create a client
	client := proto.NewMessageServiceClient(conn)

	// send message
	_, err := client.ReceiveMessage(context.Background(), message)
	if err != nil {
		return err
	}

	m.MessageLogs = append(m.MessageLogs, messageLog{
		hash:        hash(message.GetValue()),
		receiver:    conn.Target(),
		time:        time.Now(),
		messageType: int(message.Type),
	})

	return nil
}

// hash returns the hash value of data.
func hash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func (m *messageManager) ReceiveMessage(ctx context.Context, request *proto.MessageRequest) (*proto.MessageResponse, error) {

	switch request.Type {
	case proto.MessageType_QUERY:

	case proto.MessageType_DECISION:
	}
	return &proto.MessageResponse{}, nil
}
