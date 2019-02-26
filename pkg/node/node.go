package node

import (
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"go.uber.org/zap"

	//"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/network/peer"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"github.com/wavesplatform/gowaves/pkg/state"
	//"github.com/wavesplatform/gowaves/pkg/proto"
	//"github.com/wavesplatform/gowaves/pkg/state"
)

type Node struct {
	//peers Peers
	//inner        inner
	stateManager *state.StateManager
}

func NewNode(stateManager *state.StateManager) *Node {
	return &Node{
		stateManager: stateManager,
	}
}

func (a *Node) HandleProtoMessage(respondTo peer.Peer, mess peer.ProtoMessage) {
	switch t := mess.Message.(type) {
	case *proto.GetBlockMessage:
		a.handleBlockBySignature(respondTo, t.BlockID)
	}
}

func (a *Node) handleBlockBySignature(respondTo peer.Peer, sig crypto.Signature) {
	block, err := a.stateManager.GetBlock(sig)
	if err != nil {
		zap.S().Error(err)
		return
	}

	bts, err := block.MarshalBinary()
	if err != nil {
		zap.S().Error(err)
		return
	}

	bm := proto.BlockMessage{
		BlockBytes: bts,
	}

	respondTo.SendMessage(&bm)
}

//type Peers struct {
//}

//type Runtime struct {
//
//}
