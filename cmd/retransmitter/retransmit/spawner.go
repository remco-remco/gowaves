package retransmit

import (
	"context"
	"net"

	"github.com/wavesplatform/gowaves/pkg/libs/bytespool"
	"github.com/wavesplatform/gowaves/pkg/network/conn"
	"github.com/wavesplatform/gowaves/pkg/network/peer"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

type PeerSpawner interface {
	SpawnOutgoing(ctx context.Context, address string)
	SpawnIncoming(ctx context.Context, c net.Conn)
}

type PeerSpawnerImpl struct {
	pool         bytespool.Pool
	parent       peer.Parent
	wavesNetwork string
	declAddr     proto.PeerInfo
	skipFunc     conn.SkipFilter
}

func NewPeerSpawner(pool bytespool.Pool, skipFunc conn.SkipFilter, parent peer.Parent, WavesNetwork string, declAddr proto.PeerInfo) *PeerSpawnerImpl {
	return &PeerSpawnerImpl{
		pool:         pool,
		skipFunc:     skipFunc,
		parent:       parent,
		wavesNetwork: WavesNetwork,
		declAddr:     declAddr,
	}
}

func (a *PeerSpawnerImpl) SpawnOutgoing(ctx context.Context, address string) {
	params := peer.OutgoingPeerParams{
		Address:      address,
		WavesNetwork: a.wavesNetwork,
		Parent:       a.parent,
		Skip:         a.skipFunc,
		Pool:         a.pool,
		DeclAddr:     a.declAddr,
	}

	peer.RunOutgoingPeer(ctx, params)
}

func (a *PeerSpawnerImpl) SpawnIncoming(ctx context.Context, c net.Conn) {
	params := peer.IncomingPeerParams{
		WavesNetwork: a.wavesNetwork,
		Conn:         c,
		Skip:         a.skipFunc,
		Parent:       a.parent,
		DeclAddr:     a.declAddr,
		Pool:         a.pool,
	}

	peer.RunIncomingPeer(ctx, params)
}
