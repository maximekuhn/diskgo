package discovery

import (
	"context"
)

type Advertiser interface {
	// Advertise to the local network
	Advertise(ctx context.Context) error
}
