package cli

import "net"

type CmdType int

const (
	CmdShowHelp CmdType = iota
	CmdExit
	CmdAddPeer
	CmdSaveFile
	CmdList
	CmdGetFile
	CmdListPeers
)

type Command struct {
	CmdType CmdType
	Payload interface{}
}

type PayloadAddPeer struct {
	Name string
	Addr net.IP
	Port uint16
}

type PayloadSaveFile struct {
	// path of the file to save
	Path string
}

type PayloadGetFile struct {
	// name of the file to retrieve
	FileName string
}
