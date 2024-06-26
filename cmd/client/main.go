package main

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/maximekuhn/diskgo/cmd/client/cli"
	"github.com/maximekuhn/diskgo/internal/client"
	"github.com/maximekuhn/diskgo/internal/encryption"
	"github.com/maximekuhn/diskgo/internal/network"
	"github.com/maximekuhn/diskgo/internal/network/discovery"
)

func main() {
	fmt.Println(banner())

	replicas, err := client.NewBasicReplicationManager(2)
	if err != nil {
		panic(err)
	}

	c := client.NewClient(
		client.WithNickName("maxime"),
		client.WithFileEncrypter(encryption.NewAESFileEncryptor([]byte("i5yrqDhVmvV9YpFBwexikVXYFtC4emd9"))),
		client.WithResolver(discovery.NewZeroconfResolver()),
		client.WithStateStoragePath("./files/client_state.json"),
		client.WithReplicasStrategy(replicas),
	)

	// try to restore previous state, if any
	err = c.Restore()
	if err != nil {
		fmt.Printf("could not restore previous state: %s\n", err)
	} else {
		fmt.Println("successfully restored previous state")
	}

	// start peers discovery
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = c.StartDiscovery(ctx)
	if err != nil {
		panic(err)
	}

	inputCh := make(chan string, 1)
	go cli.ReadFromStdin(inputCh)

	for {
		fmt.Print("$ ")
		input := <-inputCh
		if input == "" {
			continue
		}

		cmd, err := cli.ParseCommand(input)
		if err != nil {
			fmt.Printf("failed to parse command: %s\n", err)
			continue
		}

		if cmd.CmdType == cli.CmdShowHelp {
			fmt.Println(help())
		}

		if cmd.CmdType == cli.CmdExit {
			fmt.Println("👋 bye")
			return
		}

		if cmd.CmdType == cli.CmdSaveFile {
			payload := cmd.Payload.(cli.PayloadSaveFile)
			if err := c.SaveFile(payload.Path); err != nil {
				fmt.Printf("failed to save file: %s\n", err)
				continue
			}
			fmt.Println("saved file successfully")
		}

		if cmd.CmdType == cli.CmdAddPeer {
			payload := cmd.Payload.(cli.PayloadAddPeer)
			peerAddr, err := netip.ParseAddrPort(fmt.Sprintf("%s:%d", payload.Addr, payload.Port))
			if err != nil {
				fmt.Printf("provided peer address is not valid: %s\n", err)
				continue
			}

			if err := c.AddPeer(network.NewPeer(payload.Name, peerAddr)); err != nil {
				fmt.Printf("failed to add peer: %s\n", err)
				continue
			}
			fmt.Println("added peer successfully")
		}

		if cmd.CmdType == cli.CmdGetFile {
			payload := cmd.Payload.(cli.PayloadGetFile)
			err := c.GetFile(payload.FileName)
			if err != nil {
				fmt.Printf("failed to get file: %s\n", err)
				continue
			}
			fmt.Println("got file successfully")
		}

		if cmd.CmdType == cli.CmdList {
			files := c.ListFiles()
			for filename, peersNames := range files {
				fmt.Printf("%s -> %v\n", filename, peersNames)
			}
		}

		if cmd.CmdType == cli.CmdListPeers {
			peers := c.ListPeers()
			for _, peer := range peers {
				fmt.Printf("%s -> %s\n", peer.Name, peer.Addr)
			}
		}

	}

}

func banner() string {
	return `
      _  ___  __            
 | \  |  (_  |/  _   _  
 |_/ _|_ __) |\ (_| (_)  (client)
                 _|
    type ` + "`" + `help` + "`" + ` to get the list of all available commands
    `
}

func help() string {
	return `
        * save <file path>       - save a file (a random peer will be chosen)
        * get <file name>        - retrieve a file from the peers network
        * ls | list              - list all files saved in the peers network
        * peers                  - list all known peers
        * add <name> <IP> <port> - manually add a new peer
        * help                   - show this menu
        * quit | exit            - exit
    `
}
