//go:build integration

package discovery

import (
	"context"
	"github.com/maximekuhn/diskgo/internal/network"
	"net/netip"
	"testing"
	"time"
)

func TestZeroConfBasic(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	peers := make(chan network.Peer)
	resolver := NewZeroconfResolver()
	go func() {
		if err := resolver.Resolve(ctx, peers); err != nil {
			t.Errorf("failed to resolve zeroconf: %v", err)
		}

		close(peers)
	}()

	addr, err := netip.ParseAddrPort("192.168.1.18:8888")
	if err != nil {
		t.Errorf("failed to parse addr: %v", err)
	}
	nickname := "clem"
	advertiser := NewZeroConfAdvertiser(nickname, addr)
	go func() {
		if err := advertiser.Advertise(ctx); err != nil {
			t.Errorf("failed to advertise: %v", err)
		}
	}()

	<-ctx.Done()

	// we should have discovered the server
	discoveredPeers := make([]network.Peer, 0)
	for peer := range peers {
		discoveredPeers = append(discoveredPeers, peer)
	}

	if len(discoveredPeers) < 1 {
		t.Errorf("should have at least one discovered peer")
	}

	discoveredPeer := discoveredPeers[0]
	actualName := discoveredPeer.Name
	if nickname != actualName {
		t.Errorf("wrong name got %s want %s", actualName, nickname)
	}

	actualAddr := discoveredPeer.Addr
	if addr != actualAddr {
		t.Errorf("wrong addr got %s want %s", actualAddr, actualName)
	}
}
