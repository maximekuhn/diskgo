package main

import (
	"fmt"
	"net/netip"

	"github.com/maximekuhn/diskgo/internal/client"
	"github.com/maximekuhn/diskgo/internal/file"
	"github.com/maximekuhn/diskgo/internal/network"
)

func main() {
	fmt.Println(banner())

	c := client.NewClient()

	addrPort, err := netip.ParseAddrPort("127.0.0.1:9999")
	if err != nil {
		panic(err)
	}

	filename := "readme.md"

	c.AddPeer(network.NewPeer("toto", addrPort), filename)

	err = c.GetFile(filename)
	if err != nil {
		panic(err)
	}

	err = c.SaveFile(file.NewFile(filename, []byte{1, 2, 3}))
	if err != nil {
		panic(err)
	}
}

func banner() string {
	return `
      _  ___  __            
 | \  |  (_  |/  _   _  
 |_/ _|_ __) |\ (_| (_)  (client)
                 _|
    `
}
