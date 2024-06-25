package client

import (
	"context"
	"errors"
	"fmt"
	"net/netip"

	"github.com/maximekuhn/diskgo/internal/encryption"
	"github.com/maximekuhn/diskgo/internal/network/discovery"

	"github.com/maximekuhn/diskgo/internal/file"
	"github.com/maximekuhn/diskgo/internal/network"
	"github.com/maximekuhn/diskgo/internal/protocol"
)

type Client struct {
	manager          *peersManager
	replicasStrategy ReplicasStrategy

	// can be nil
	fileEncrypter encryption.FileEncrypter

	// can be nil
	resolver discovery.Resolver

	nickname         string
	stateStoragePath string
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

func (c *Client) StartDiscovery(ctx context.Context) error {
	peers := make(chan network.Peer)
	if c.resolver != nil {
		if err := c.resolver.Resolve(ctx, peers); err != nil {
			return err
		}

		go func(ctx context.Context, c *Client) {
			for {
				select {
				case peer := <-peers:
					// error only indicates a duplicate
					if err := c.AddPeer(&peer); err == nil {
						fmt.Printf("discovered %s at %s\n", peer.Name, peer.Addr.String())
					}

				case <-ctx.Done():
					return
				}
			}
		}(ctx, c)
	}

	return nil
}

func (c *Client) SaveFile(filepath string) error {
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

	// save the file according to the provided replication strategy
	err = c.replicasStrategy.Save(f, c.manager, c.nickname)
	if err != nil {
		return err
	}

	// snapshot client's state and save it (to disk)
	if c.stateStoragePath != "" {
		err = c.snapshot()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) AddPeer(p *network.Peer) error {
	return c.manager.addPeer(p)
}

func (c *Client) GetFile(filename string) error {
	peer, err := c.manager.getPeerStoringFile(filename)
	if err != nil {
		return err
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
	res, err := sendRequest(context.TODO(), &req, peer)
	if err != nil {
		return err
	}

	if res.MsgType != protocol.MsgGetFileRes {
		return errors.New("got a response but not the one expected")
	}

	payload := res.Payload.(protocol.GetFileResPayload)
	if !payload.Ok {
		return errors.New("failed to get file from peer")
	}

	f := &payload.File

	// check if we need to decrypt the file
	if c.fileEncrypter != nil {
		err = c.fileEncrypter.Decrypt(f)
		if err != nil {
			return err
		}
	}

	// save the file on disk
	if err = file.WriteFile(f); err != nil {
		return err
	}

	return nil
}

func (c *Client) ListFiles() map[string][]string {
	return c.manager.files
}

func (c *Client) ListPeers() []network.Peer {
	peers := make([]network.Peer, 0)
	for _, peer := range c.manager.knownPeers {
		peers = append(peers, *peer)
	}
	return peers
}

func (c *Client) snapshot() error {
	peers := c.ListPeers()
	statePeers := make([]statePeer, 0)
	for _, peer := range peers {
		statePeers = append(statePeers, statePeer{
			Name: peer.Name,
			Addr: peer.Addr.String(),
		})
	}

	files := c.ListFiles()
	stateFiles := make([]stateFile, 0)
	for filename, peernames := range files {
		for _, peername := range peernames {
			stateFiles = append(stateFiles, stateFile{
				Filename: filename,
				Peername: peername,
			})
		}
	}

	m := &clientPersistence{
		clientState: state{
			Peers: statePeers,
			Files: stateFiles,
		},
		writeFilePath: c.stateStoragePath,
	}

	return m.writeToDisk()
}

func (c *Client) Restore() error {
	m := newClientPeristence(c.stateStoragePath)
	state, err := m.readFromDisk()
	if err != nil {
		return err
	}

	peers := state.Peers
	for _, peer := range peers {
		peerName := peer.Name
		peerAdrr, err := netip.ParseAddrPort(peer.Addr)
		if err != nil {
			return err
		}

		peer := network.NewPeer(peerName, peerAdrr)
		err = c.manager.addPeer(peer)
		if err != nil {
			return err
		}
	}

	files := state.Files
	for _, f := range files {
		peerName := f.Peername
		fileName := f.Filename

		// FIXME: manager only use the peer's Name but that's not correct to do so (even if it works)
		peer := &network.Peer{
			Name: peerName,
		}

		c.manager.addFilePeerStorage(fileName, peer)
	}

	return nil
}
