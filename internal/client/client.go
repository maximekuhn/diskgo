package client

import (
	"errors"

	"github.com/maximekuhn/diskgo/internal/file"
	"github.com/maximekuhn/diskgo/internal/network"
	"github.com/maximekuhn/diskgo/internal/protocol"
)

type Client struct {
	manager *peersManager
}

func NewClient() *Client {
	return &Client{
		manager: newPeersManager(),
	}
}

func (c *Client) SaveFile(filepath string) error {
	// get a random known peer to save file
	peer, err := c.manager.getRandomPeer()
	if err != nil {
		return err
	}

	// read the file
	f, err := file.ReadFile(filepath)
	if err != nil {
		return err
	}

	// send the request
	req := protocol.Message{
		MsgType: protocol.MsgSaveFile,
		Payload: protocol.SaveFileReqPayload{
			File: *f,
		},
	}

	// await for response
	res, err := sendRequest(&req, peer)
	if err != nil {
		return err
	}

	if res.MsgType != protocol.MsgSaveFileRes {
		return errors.New("got a response but not the one expected")
	}

	// maybe check the payload cast (should be done by encoding/decoding)
	payload := res.Payload.(protocol.SaveFileResPayload)
	if !payload.Ok {
		return errors.New("peer failed to save the file")
	}

	c.manager.addFilePeerStorage(f.Name, peer)

	return nil
}

func (c *Client) AddPeer(p *network.Peer) error {
	return c.manager.addPeer(p)
}

func (c *Client) GetFile(filename string) (*file.File, error) {
	peer, err := c.manager.getPeerStoringFile(filename)
	if err != nil {
		return nil, err
	}

	// send the request
	req := protocol.Message{
		MsgType: protocol.MsgGetFile,
		Payload: protocol.GetFileReqPayload{
			FileName: filename,
		},
	}

	// await response from remote peer
	res, err := sendRequest(&req, peer)
	if err != nil {
		return nil, err
	}

	if res.MsgType != protocol.MsgGetFileRes {
		return nil, errors.New("got a response but not the one expected")
	}

	payload := res.Payload.(protocol.GetFileResPayload)
	if !payload.Ok {
		return nil, errors.New("failed to get file from peer")
	}

	return &payload.File, nil
}

func (c *Client) ListFiles() map[string][]string {
	return c.manager.files
}
