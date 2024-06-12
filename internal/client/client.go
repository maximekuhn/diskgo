package client

import (
	"errors"
	"fmt"
	"github.com/maximekuhn/diskgo/internal/encryption"

	"github.com/maximekuhn/diskgo/internal/file"
	"github.com/maximekuhn/diskgo/internal/network"
	"github.com/maximekuhn/diskgo/internal/protocol"
)

type Client struct {
	manager *peersManager

	// can be nil
	fileEncrypter encryption.FileEncrypter
	nickname      string
}

func NewClient(opts ...ClientOpts) *Client {
	// create client with default opts
	c := &Client{
		manager: newPeersManager(),
	}
	for _, opt := range DefaultClientOpts() {
		opt(c)
	}

	// apply provided opts
	for _, opt := range opts {
		opt(c)
	}

	return c
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

	// check if encryption is required
	if c.fileEncrypter != nil {
		err = c.fileEncrypter.Encrypt(f)
		if err != nil {
			return err
		}
	}

	// send the request
	req := protocol.Message{
		MsgType: protocol.MsgSaveFile,
		Payload: protocol.SaveFileReqPayload{
			File: *f,
		},
		From: c.nickname,
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
		return fmt.Errorf("peer failed to save the file '%s'", payload.Reason)
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
		From: c.nickname,
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

	f := &payload.File

	// check if we need to decrypt the file
	if c.fileEncrypter != nil {
		err = c.fileEncrypter.Decrypt(f)
		if err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (c *Client) ListFiles() map[string][]string {
	return c.manager.files
}
