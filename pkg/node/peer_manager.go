package node

import (
	"github.com/wavesplatform/gowaves/pkg/network/peer"
	"sync"
)

type PeerManager interface {
	Connected(unique string) (peer.Peer, bool)
	Banned(unique string) bool
	AddConnected(p PeerInfo)
	PeerWithHighestScore() (peer.Peer, uint64, bool)
}

type PeerManagerImpl struct {
	spawner           PeerSpawner
	activeConnections map[string]peer.Peer
	mu                sync.RWMutex
}

func (a *PeerManagerImpl) Connected(unique string) (peer.Peer, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	p, ok := a.activeConnections[unique]
	return p, ok
}
