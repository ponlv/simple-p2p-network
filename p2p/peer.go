package p2p

import (
	"google.golang.org/grpc"
	"time"
)

// peer is the remote node that a local node can connect to.
type peer struct {
	Address string           // network address
	conn    *grpc.ClientConn // client connection
}

var (
	maxPeerNum           = 20              // max neighbor peers' number
	maxDiscoverSleepTime = 5 * time.Second // sleep time between discover neighbor peers
)
