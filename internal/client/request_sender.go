package client

import (
	"fmt"
	"net"

	"github.com/maximekuhn/diskgo/internal/network"
	"github.com/maximekuhn/diskgo/internal/protocol"
)

func sendRequest(req *protocol.Message, peer *network.Peer) (*protocol.Message, error) {
	encodedMsg, err := network.Encode(req.MsgType, req.Payload)
	if err != nil {
		return nil, err
	}

	peerAddr := fmt.Sprintf("%s:%d", peer.Addr.Addr(), peer.Addr.Port())

	conn, err := net.Dial("tcp", peerAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// TODO: handle n write
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
