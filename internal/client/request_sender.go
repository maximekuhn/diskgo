package client

import (
	"context"
	"fmt"
	"net"

	"github.com/maximekuhn/diskgo/internal/network"
	"github.com/maximekuhn/diskgo/internal/protocol"
)

// TODO: handle context for write and read

func sendRequest(ctx context.Context, req *protocol.Message, peer *network.Peer) (*protocol.Message, error) {
	encodedMsg, err := network.Encode(*req)
	if err != nil {
		return nil, err
	}

	peerAddr := fmt.Sprintf("%s:%d", peer.Addr.Addr(), peer.Addr.Port())

	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", peerAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	_, err = conn.Write(encodedMsg)
	if err != nil {
		return nil, err
	}

	decodedRes, err := network.Decode(conn)
	if err != nil {
		return nil, err
	}

	return decodedRes, nil
}
