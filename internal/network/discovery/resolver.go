package discovery

import (
	"context"
	"github.com/maximekuhn/diskgo/internal/network"
)

type Resolver interface {
	// Resolve listens for advertisement on the local network
	Resolve(ctx context.Context, peers chan<- network.Peer) error
}
