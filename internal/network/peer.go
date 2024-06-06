package network

import "net/netip"

type Peer struct {
	Name string
	Addr netip.AddrPort
}

func NewPeer(name string, addr netip.AddrPort) *Peer {
	return &Peer{
		Name: name,
		Addr: addr,
	}
}
