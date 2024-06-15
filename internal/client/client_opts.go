package client

import (
	"fmt"
	"github.com/maximekuhn/diskgo/internal/encryption"
	"github.com/maximekuhn/diskgo/internal/network/discovery"
	"math/rand"
)

type ClientOpts func(*Client)

func WithFileEncrypter(encrypter encryption.FileEncrypter) ClientOpts {
	return func(c *Client) {
		c.fileEncrypter = encrypter
	}
}

func WithNickName(nickname string) ClientOpts {
	return func(c *Client) {
		c.nickname = nickname
	}
}

func WithResolver(resolver discovery.Resolver) ClientOpts {
	return func(c *Client) {
		c.resolver = resolver
	}
}

// WithStateStoragePath determines where the client's state is saved to be restored across client restart(s)
//
// Set it to an empty for non-persistent state
func WithStateStoragePath(stateStoragePath string) ClientOpts {
	return func(c *Client) {
		c.stateStoragePath = stateStoragePath
	}
}

func DefaultClientOpts() []ClientOpts {
	// random nickname, not recommended
	randomNickname := fmt.Sprintf("user-%d", rand.Intn(100))
	return []ClientOpts{
		WithFileEncrypter(nil),       // no encryption
		WithNickName(randomNickname), // use a random nickname
		WithResolver(nil),            // not automatic resolver
		WithStateStoragePath(""),     // non-persistent state
	}
}
