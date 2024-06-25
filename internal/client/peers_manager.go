package client

import (
	"errors"
	"math/rand"

	"github.com/maximekuhn/diskgo/internal/network"
)

// XXX: maybe make this thread safe

type peersManager struct {
	// peer Name -> peer
	knownPeers map[string]*network.Peer

	// Filename -> peer Name
	files map[string][]string
}

func newPeersManager() *peersManager {
	return &peersManager{
		knownPeers: make(map[string]*network.Peer),
		files:      make(map[string][]string),
	}
}

// add the given peer
//
// if a peer with the same Name already exists, an error is returned and the peer isn't added
func (m *peersManager) addPeer(peer *network.Peer) error {
	p, ok := m.knownPeers[peer.Name]
	if !ok {
		m.knownPeers[peer.Name] = peer
		return nil
	}

	if p != nil {
		return errors.New("peer with same Name already exists")
	}

	m.knownPeers[peer.Name] = peer
	return nil
}

// add a new peer as a storage for the given Filename
func (m *peersManager) addFilePeerStorage(filename string, peer *network.Peer) {
	ps, ok := m.files[filename]
	if !ok {
		ps = make([]string, 0)
	}

	ps = append(ps, peer.Name)
	m.files[filename] = ps
}

// get a peer storing the given Filename
//
// if there is no peer, an error is returned
func (m *peersManager) getPeerStoringFile(filename string) (*network.Peer, error) {
	peersNames, ok := m.files[filename]
	if !ok {
		return nil, errors.New("no peer is storing this file")
	}

	peersCount := len(peersNames)
	if peersCount == 0 {
		return nil, errors.New("no peer is storing this file")
	}

	idx := rand.Intn(peersCount)
	peerName := peersNames[idx]

	peer, ok := m.knownPeers[peerName]
	if !ok {
		return nil, errors.New("a peer has this file but I don't know him")
	}

	return peer, nil
}

// getAllKnownPeers return the list of all currently known peers
func (m *peersManager) getAllKnownPeers() []*network.Peer {
	peers := make([]*network.Peer, 0)
	for _, p := range m.knownPeers {
		peers = append(peers, p)
	}
	return peers
}
