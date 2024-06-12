package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/maximekuhn/diskgo/internal/server"
	"github.com/maximekuhn/diskgo/internal/store"
)

func main() {
	// parse CLI args
	port := flag.Int("port", 9999, "port to listen on")
	flag.Parse()

	// show config
	fmt.Println(banner())
	fmt.Println("port", *port)
	fmt.Println()

	// start the server
	s := server.NewServer(
		server.WithListenPort(uint16(*port)),
		server.WithListenAddr(net.IPv4(0, 0, 0, 0)),
		server.WithFileStore(store.NewFsFileStore("./files", 0)),
	)

	stopCh := make(chan bool, 1)
	err := s.Start(stopCh)
	if err != nil {
		slog.Error("failed to start server", slog.String("err_msg", err.Error()))
		os.Exit(1)
	}

	// listen on signal to stop the server
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGABRT)
	<-signalCh

	// stop the server
	stopCh <- true
}

func banner() string {
	return `
      _  ___  __            
 | \  |  (_  |/  _   _  
 |_/ _|_ __) |\ (_| (_)  (server)
                 _|
    `
}
