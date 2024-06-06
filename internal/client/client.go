package client

import (
	"errors"
	"fmt"
	"net"

	"github.com/maximekuhn/diskgo/internal/file"
	"github.com/maximekuhn/diskgo/internal/network"
	"github.com/maximekuhn/diskgo/internal/protocol"
)

type Client struct {
	// file name -> peers
	myFiles map[string][]*network.Peer
	peers   []*network.Peer
}

func NewClient() *Client {
	return &Client{
		myFiles: make(map[string][]*network.Peer),
		peers:   make([]*network.Peer, 0),
	}
}

// TODO: remove this, it's only for early dev purposes
func (c *Client) AddPeer(p *network.Peer, filename string) {
	fp, ok := c.myFiles[filename]
	if !ok {
		c.myFiles[filename] = make([]*network.Peer, 0)
		fp = c.myFiles[filename]
	}

	fp = append(fp, p)
	c.myFiles[filename] = fp
	c.peers = append(c.peers, p)
}

func (c *Client) GetFile(filename string) error {
	peers, ok := c.myFiles[filename]
	if !ok {
		return errors.New("no peers has this file")
	}
	if len(peers) == 0 {
		return errors.New("no peers has this file")
	}

	// TODO: ask a random peer
	peer := peers[0]

	payload := protocol.GetFileReqPayload{
		FileName: filename,
	}

	encodedMsg, err := network.Encode(protocol.MsgGetFile, payload)
	if err != nil {
		return err
	}

	// dial the remote peer and send the request
	addr := fmt.Sprintf("%s:%d", peer.Addr.Addr(), peer.Addr.Port())
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// TODO: handle n write
	_, err = conn.Write(encodedMsg)
	if err != nil {
		return err
	}

	// wait for remote peer to give a response
	decodedMsg, err := network.Decode(conn)
	if err != nil {
		return err
	}

	if decodedMsg.MsgType != protocol.MsgGetFileRes {
		return errors.New("received a response but not the one expected")
	}

	res, ok := decodedMsg.Payload.(protocol.GetFileResPayload)
	if !ok {
		return errors.New("invalid payload")
	}

	if !res.Ok {
		fmt.Println("the peer doesn't have the file :(")
	} else {
		fmt.Printf("received file %s\ndata: %v\n", res.File.Name, res.File.Data)
	}

	return nil
}

func (c *Client) SaveFile(f *file.File) error {
	// TODO: choose a random peer
	if len(c.peers) == 0 {
		return errors.New("no peers known")
	}
	peer := c.peers[0]

	// create request
	req := protocol.SaveFileReqPayload{
		File: *f,
	}
	encodedMsg, err := network.Encode(protocol.MsgSaveFile, req)
	if err != nil {
		return err
	}

	// send request
	addr := fmt.Sprintf("%s:%d", peer.Addr.Addr(), peer.Addr.Port())
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	// TODO: handle n write
	_, err = conn.Write(encodedMsg)
	if err != nil {
		return err
	}

	// wati for remote peer to respond
	decodedMsg, err := network.Decode(conn)
	if err != nil {
		return err
	}

	if decodedMsg.MsgType != protocol.MsgSaveFileRes {
		return errors.New("received a response but not the one expected")
	}
	res, ok := decodedMsg.Payload.(protocol.SaveFileResPayload)
	if !ok {
		return errors.New("invalid payload")
	}

	if res.Ok {
		fmt.Println("peer successfully saved file")
	} else {
		fmt.Println("peer could not save file")
	}

	return nil
}
