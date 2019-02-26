package node

import (
	"github.com/stretchr/testify/require"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/network/conn"
	"github.com/wavesplatform/gowaves/pkg/network/peer"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"github.com/wavesplatform/gowaves/pkg/state"
	"github.com/wavesplatform/gowaves/pkg/util"
	"testing"
)

type mockPeer struct {
}

func (*mockPeer) Direction() peer.Direction {
	panic("implement me")
}

func (*mockPeer) Close() {
	panic("implement me")
}

func (*mockPeer) SendMessage(proto.Message) {
	panic("implement me")
}

func (*mockPeer) ID() string {
	panic("implement me")
}

func (*mockPeer) Connection() conn.Connection {
	panic("implement me")
}

func TestNode_HandleProtoMessage_GetBlockBySignature(t *testing.T) {
	dataDir, _ := util.NewTemporary()
	defer dataDir.Clear()

	s, err := state.NewStateManager(string(dataDir), state.DefaultBlockStorageParams())
	require.NoError(t, err)

	n := NewNode(s)

	sig, _ := crypto.NewSignatureFromBase58("5uqnLK3Z9eiot6FyYBfwUnbyid3abicQbAZjz38GQ1Q8XigQMxTK4C1zNkqS1SVw7FqSidbZKxWAKLVoEsp4nNqa")

	n.handleBlockBySignature(&mockPeer{}, sig)

}
