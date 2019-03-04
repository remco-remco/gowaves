package node

import (
	"context"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/network/peer"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"go.uber.org/zap"
	"net"
	"time"
)

type Config struct {
	AppName  string
	NodeName string
	Listen   string
	DeclAddr string
}

type StateManager interface {
	GetBlock(blockID crypto.Signature) (*proto.Block, error)
	AddBlocks(blocks [][]byte, initialisation bool) error
	AddBlock(block *proto.Block) error
}

type Node struct {
	peerManager  PeerManager
	stateManager StateManager
	subscribe    *Subscribe
}

func NewNode(stateManager StateManager, peerManager PeerManager) *Node {
	return &Node{
		stateManager: stateManager,
		peerManager:  peerManager,
	}
}

func (a *Node) HandleProtoMessage(respondTo string, mess peer.ProtoMessage) {
	switch t := mess.Message.(type) {
	case *proto.GetBlockMessage:
		a.handleBlockBySignature(respondTo, t.BlockID)
	}
}

func (a *Node) HandleInfoMessage(m peer.InfoMessage) {
	switch t := m.Value.(type) {
	case *peer.Connected:
		a.handleConnected(m.ID, t)
	}
}

func (a *Node) handleConnected(id string, t *peer.Connected) {
	peerInfo := PeerInfo{
		Peer: t.Peer,
	}

	_, connected := a.peerManager.Connected(peerInfo.Unique())
	if connected {
		peerInfo.Peer.Close()
		return
	}

	if a.peerManager.Banned(peerInfo.Unique()) {
		peerInfo.Peer.Close()
		return
	}

	a.peerManager.AddConnected(peerInfo)
}

func (a *Node) handleBlockBySignature(peer string, sig crypto.Signature) {
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

	p, ok := a.peerManager.Connected(peer)
	if ok {
		p.SendMessage(&bm)
	}
}

// called every n seconds, handle change runtime state
func (a *Node) Tick() {

	for {
		p, score, ok := a.peerManager.PeerWithHighestScore()
		if !ok {
			// no peers, skip
			return
		}

		if score == 0 {
			time.Sleep(5 * time.Second)
			continue
		}

		// TODO check if we have highest score

		p.SendMessage(&proto.GetSignaturesMessage{})

		messCh, unsubscribe := a.subscribe.Subscribe(p, &proto.SignaturesMessage{})

		var mess *proto.SignaturesMessage

		select {
		case <-time.After(15 * time.Second):
		// TODO handle timeout
		case received := <-messCh:
			//a.subscribe.Unsubscribe(p, &proto.SignaturesMessage{})
			unsubscribe()
			mess = received.(*proto.SignaturesMessage)
		}

		blockSignatures := BlockSignatures{}

		funcName(mess, blockSignatures, p, a)

		//?, ? := a.findMaxCommonBlock(mess.Signatures)

		//for _, i := range mess.Signatures {
		//}

		//if err != nil {
		//	if err == TimeoutErr {
		//		// TODO handle timeout
		//	}
		//}

		//ask.Subscribe(15*time.Second)
		//
		//a.subscribe.Clear(ask)
		//
		//if ask.Timeout() {
		//	// TODO handle timeout
		//}
		//
		//m := ask.Get().(*proto.SignaturesMessage{})

	}

}

func funcName(mess *proto.SignaturesMessage, blockSignatures BlockSignatures, p peer.Peer, a *Node) {
	subscribeCh, unsubscribe := a.subscribe.Subscribe(p, &proto.BlockMessage{})
	defer unsubscribe()
	for _, sig := range mess.Signatures {
		if !blockSignatures.Exists(sig) {
			p.SendMessage(&proto.GetBlockMessage{BlockID: sig})

			// wait for block with expected signature
			timeout := time.After(30 * time.Second)
			for {
				select {
				case <-timeout:
				// TODO HANDLE timeout

				case blockMessage := <-subscribeCh:
					block := proto.Block{}
					err := block.UnmarshalBinary(blockMessage.(*proto.BlockMessage).BlockBytes)
					if err != nil {
						zap.S().Error(err)
						continue
					}

					if block.BlockSignature != sig {
						continue
					}

					err = a.stateManager.AddBlock(&block)
					if err != nil {
						// TODO handle error
					}
					break
				}
			}
		}
	}
}

func RunIncomeConnectionsServer(ctx context.Context, n *Node, c Config, s PeerSpawner) {
	l, err := net.Listen("tcp", c.Listen)
	if err != nil {
		zap.S().Error(err)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			zap.S().Error(err)
			continue
		}

		go s.SpawnIncoming(ctx, c)
	}
}

func RunNode(n *Node) {

}

type BlockSignatures struct {
	signatures []crypto.Signature
	unique     map[crypto.Signature]struct{}
}

func (a *BlockSignatures) Exists(sig crypto.Signature) bool {
	_, ok := a.unique[sig]
	return ok
}

//type Peers struct {
//}

//type Runtime struct {
//
//}
