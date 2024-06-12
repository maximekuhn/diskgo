package server

import (
	"context"
	"fmt"
	"github.com/maximekuhn/diskgo/internal/network"
	"github.com/maximekuhn/diskgo/internal/network/discovery"
	"github.com/maximekuhn/diskgo/internal/protocol"
	"github.com/maximekuhn/diskgo/internal/server/handlers"
	"github.com/maximekuhn/diskgo/internal/store"
	"log/slog"
	"net"
)

type Server struct {
	addr  net.IP
	port  uint16
	store store.FileStore

	// can be nil
	advertiser discovery.Advertiser
}

func NewServer(opts ...ServerOpts) *Server {
	// create server with default opts
	s := &Server{}
	for _, opt := range DefaultServerOpts() {
		opt(s)
	}

	// apply provided opts
	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Start(stopCh <-chan bool) error {
	listenAddr := fmt.Sprintf("%s:%d", s.addr, s.port)
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	slog.Info("server started", slog.String("listen_addr", l.Addr().String()))

	if s.advertiser != nil {
		ctx := context.Background()

		// TODO: handle error
		go s.advertiser.Advertise(ctx)
	}

	s.mainLoop(l, stopCh)

	return nil
}

func (s *Server) mainLoop(l net.Listener, stopCh <-chan bool) {
	defer l.Close()

	connCh := make(chan net.Conn)
	go acceptConnLoop(l, connCh)

	for {
		select {
		case conn := <-connCh:
			slog.Info("accepted incoming conn", slog.String("remote_addr", conn.RemoteAddr().String()))
			go s.handleConn(conn)

		case <-stopCh:
			slog.Info("received stop signal")
			return
		}
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	msg, err := network.Decode(conn)

	if err != nil {
		slog.Error("failed to decode incoming message", slog.String("err_msg", err.Error()))
		return
	}

	slog.Info("received a message", slog.String("nickname", msg.From))

	if msg.MsgType == protocol.MsgGetFile {
		slog.Info("received a MsgGetFile request")

		err := handlers.HandleGetFile(msg, s.store, conn)
		if err != nil {
			slog.Error("failed to handle request", slog.String("err_msg", err.Error()))
			return
		}
	}

	if msg.MsgType == protocol.MsgSaveFile {
		slog.Info("received a MsgSaveFile request")

		err := handlers.HandleSaveFile(msg, s.store, conn)
		if err != nil {
			slog.Error("failed to handle request", slog.String("err_msg", err.Error()))
			return
		}
	}

	slog.Info("successfully handled request")
}

func acceptConnLoop(l net.Listener, connCh chan<- net.Conn) {
	for {
		conn, err := l.Accept()
		if err != nil {
			slog.Error("failed to accept incoming conn", slog.String("err_msg", err.Error()))
			continue
		}

		connCh <- conn
	}
}
