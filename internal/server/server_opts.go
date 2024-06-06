package server

import (
	"net"

	"github.com/maximekuhn/diskgo/internal/store"
)

type ServerOpts func(*Server)

func WithListenAddr(addr net.IP) ServerOpts {
	return func(s *Server) {
		s.addr = addr
	}
}

func WithListenPort(port uint16) ServerOpts {
	return func(s *Server) {
		s.port = port
	}
}

func WithFileStore(store store.FileStore) ServerOpts {
	return func(s *Server) {
		s.store = store
	}
}

func DefaultServerOpts() []ServerOpts {
	return []ServerOpts{
		WithListenAddr(net.IPv4(127, 0, 0, 1)),
		WithListenPort(9292),
		WithFileStore(store.NewInMemoryFileStore()),
	}
}
