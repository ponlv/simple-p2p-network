package p2p

import (
	"fmt"
	"sync"
)

type PeerManager interface {
	// AddPeers add list of addresses to the peer manager.
	AddPeers(addrs []string)

	// RemovePeer removes a peer from the peer manager.
	RemovePeer(addr string) error

	// RemoveAllPeers removes all peers from the peer manager.
	RemoveAllPeers() error

	// Disconnect closes the connection to the peer.
	Disconnect(addr string) error
}

var _ PeerManager = (*peerManager)(nil)

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
func NewPeerManager(add string) PeerManager {
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
func (pm *peerManager) AddPeers(addrs []string) {
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