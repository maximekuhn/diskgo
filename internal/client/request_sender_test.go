package client

import (
	"context"
	"net/netip"
	"testing"
	"time"

	"github.com/maximekuhn/diskgo/internal/network"
	"github.com/maximekuhn/diskgo/internal/protocol"
)

// check that send request stops when the timeout exceeds
func Test_sendRequestTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// remote addr that should not be dialable and therefore timeout
	addr, err := netip.ParseAddrPort("127.0.0.1:54673")
	if err != nil {
		panic(err)
	}
	peer := network.NewPeer("peer-1", addr)

	req := protocol.Message{
		MsgType: protocol.MsgGetFile,
		From:    "me",
		Payload: nil,
	}

	_, err = sendRequest(ctx, &req, peer)
	if err == nil {
		t.Errorf("sendRequest should have return an error")
	}
}
