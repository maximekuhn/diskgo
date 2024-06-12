package main

import (
	"flag"
	"fmt"
	"github.com/maximekuhn/diskgo/internal/network/discovery"
	"log/slog"
	"math/rand"
	"net"
	"net/netip"
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

	// create the server
	nickname := fmt.Sprintf("server-%d", rand.Intn(1000))
	addr, err := getLanAddr()
	if err != nil {
		panic(err)
	}
	addrPort, err := netip.ParseAddrPort(fmt.Sprintf("%s:%d", addr, *port))
	if err != nil {
		panic(err)
	}

	s := server.NewServer(
		server.WithListenPort(uint16(*port)),
		server.WithListenAddr(net.IPv4(0, 0, 0, 0)),
		server.WithFileStore(store.NewFsFileStore("./files", 1)),
		server.WithAdvertiser(discovery.NewZeroConfAdvertiser(nickname, addrPort)),
	)

	// start the server
	stopCh := make(chan bool, 1)
	err = s.Start(stopCh)
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

func getLanAddr() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP, nil
				}
			}
		}
	}
	return nil, nil
}

func banner() string {
	return `
      _  ___  __            
 | \  |  (_  |/  _   _  
 |_/ _|_ __) |\ (_| (_)  (server)
                 _|
    `
}
