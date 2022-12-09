package p2p

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"log"
	"simple-p2p/proto/proto"
	"sync"
	"time"
)

type Peer interface {
	// AddPeers add list of addresses to the peer manager.
	AddPeers(addrs ...string)

	// RemovePeer removes a peer from the peer manager.
	RemovePeer(addr string) error

	// RemoveAllPeers removes all peers from the peer manager.
	RemoveAllPeers() error

	// Disconnect closes the connection to the peer.
	Disconnect(addr string) error

	// GetPeers returns all the peers' addresses in the peer manager.
	GetPeers() []string

	// GetPeerState returns the state of connection to a peer
	GetPeerState(addr string) connectivity.State

	// GetConnection returns a connection to a peer.
	GetConnection(addr string) (*grpc.ClientConn, error)

	// GetPeersNum returns the number of peers in the peer manager.
	GetPeersNum() int

	// GetNeighbours returns the neighbours of a peer.
	GetNeighbours(ctx context.Context, req *proto.GetNeighbourRequest) (*proto.GetNeighbourResponse, error)

	// StartDiscoverPeers starts the peer discovery process.
	StartDiscoverPeers(bootstraps ...string)
}

// peer is the remote node that a local node can connect to.
type peer struct {
	Address string           // network address
	conn    *grpc.ClientConn // client connection
}

var (
	maxPeerNum           = 20              // max neighbor peers' number
	maxDiscoverSleepTime = 5 * time.Second // sleep time between discover neighbor peers
)

var _ Peer = (*peerManager)(nil)

// PeerManager manages the peers that a local node known.
type peerManager struct {
	addr string // network address of local node

	Peers map[string]*peer // known remote peers
	Mux   sync.RWMutex     // mutual exclusion lock for peers

	stopDiscover    chan struct{}  // stop discover neighbor peers signal
	discoverStopped chan struct{}  // discover neighbor peers stopped signal
	waiter          sync.WaitGroup // wait background goroutines
}

// NewPeerManager returns a new peer manager with its own network address.
func NewPeerManager(add string) Peer {
	return &peerManager{
		addr:            add,
		Peers:           make(map[string]*peer),
		Mux:             sync.RWMutex{},
		stopDiscover:    make(chan struct{}),
		discoverStopped: make(chan struct{}),
		waiter:          sync.WaitGroup{},
	}
}

// addPeer adds an address to the peer manager.
func (pm *peerManager) addPeer(addr string) {
	pm.Mux.Lock()
	defer pm.Mux.Unlock()

	if _, ok := pm.Peers[addr]; ok {
		return
	}

	pm.Peers[addr] = &peer{Address: addr}
}

// AddPeers add list of addresses to the peer manager.
func (pm *peerManager) AddPeers(addrs ...string) {
	for _, addr := range addrs {
		pm.addPeer(addr)
	}
}

// RemovePeer removes a peer from the peer manager. It disconnects the connection relative to the peer before removing.
func (pm *peerManager) RemovePeer(addr string) error {
	pm.Mux.Lock()
	defer pm.Mux.Unlock()

	if _, ok := pm.Peers[addr]; ok {
		if err := pm.disconnect(addr); err != err {
			return err
		}

		delete(pm.Peers, addr)
	}
	return nil
}

// RemoveAllPeers removes all peers from the peer manager.
func (pm *peerManager) RemoveAllPeers() error {

	for addr := range pm.Peers {
		if err := pm.RemovePeer(addr); err != nil {
			return err
		}
	}
	return nil
}

// Disconnect closes the connection to the peer.
func (pm *peerManager) Disconnect(addr string) error {
	pm.Mux.Lock()
	defer pm.Mux.Unlock()

	return pm.disconnect(addr)
}

// disconnect closes the connection to the peer
func (pm *peerManager) disconnect(addr string) error {
	p, ok := pm.Peers[addr]
	if !ok {
		return fmt.Errorf("%v failed to disconnect: unknown peer: %v", pm.addr, addr)
	}

	if p.conn != nil {
		if err := p.conn.Close(); err != nil {
			return fmt.Errorf("%v failed to disconnect: %v", pm.addr, err)
		}
	}
	return nil
}

// GetPeers returns all the peers' addresses in the peer manager.
func (pm *peerManager) GetPeers() []string {
	pm.Mux.RLock()
	defer pm.Mux.RUnlock()

	var addresses []string
	for _, p := range pm.Peers {
		addresses = append(addresses, p.Address)
	}
	return addresses
}

// GetPeerState returns the state of connection to a peer
func (pm *peerManager) GetPeerState(addr string) connectivity.State {
	pm.Mux.RLock()
	defer pm.Mux.RUnlock()

	p, ok := pm.Peers[addr]
	if !ok || p.conn == nil {
		return connectivity.State(-1)
	}

	return p.conn.GetState()
}

// GetConnection returns a connection to a peer.
func (pm *peerManager) GetConnection(addr string) (*grpc.ClientConn, error) {
	pm.Mux.Lock()
	defer pm.Mux.Unlock()

	p, ok := pm.Peers[addr]
	if !ok {
		pm.addPeer(addr)
	}

	p, ok = pm.Peers[addr]
	if !ok {
		return nil, fmt.Errorf("failed to get connection to peer: %v", addr)
	}

	if p.conn == nil || p.conn.GetState() == connectivity.Shutdown {
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		p.conn = conn
	}
	return p.conn, nil
}

// Wait keeps the peer manager running in background.
func (pm *peerManager) Wait() {
	pm.waiter.Wait()
}

// discoverPeers discovers new peers from another peer, and add new peers into known peers list.
func (pm *peerManager) discoverPeers(addr string) {
	conn, err := pm.GetConnection(addr)
	if err != nil {
		return
	}

	client := proto.NewPeerServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	peers, err := client.GetNeighbours(ctx, &proto.GetNeighbourRequest{Address: pm.addr})
	if err != nil {
		log.Printf("%v failed to get neighbors of peer: %v: %v", pm.addr, addr, err)
		return
	}

	pm.AddPeers(peers.Peers...)
}

// GetPeersNum returns the number of peers in the peer manager.
func (pm *peerManager) GetPeersNum() int {
	pm.Mux.RLock()
	defer pm.Mux.RUnlock()

	return len(pm.Peers)
}

// StartDiscoverPeers starts discovering new peers via bootstraps.
func (pm *peerManager) StartDiscoverPeers(bootstraps ...string) {
	pm.AddPeers(bootstraps...)

	pm.waiter.Add(1)
	go func() {
		for {
			if pm.GetPeersNum() < maxPeerNum {
				for _, addr := range pm.GetPeers() {
					pm.discoverPeers(addr)
					if pm.GetPeersNum() >= maxPeerNum {
						break
					}
				}
			}

			select {
			case <-pm.stopDiscover:
				pm.waiter.Done()
				pm.discoverStopped <- struct{}{}
				return
			case <-time.After(maxDiscoverSleepTime):
				continue
			}
		}
	}()
}

// GetNeighbors returns the already known neighbor peers, and add the requester into the known,
// peers list if it's not known before.
func (pm *peerManager) GetNeighbours(ctx context.Context, req *proto.GetNeighbourRequest) (*proto.GetNeighbourResponse, error) {
	addresses := pm.GetPeers()
	pm.AddPeers(req.Address)
	return &proto.GetNeighbourResponse{Peers: addresses}, nil
}
