package cli

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func ParseCommand(input string) (*Command, error) {
	if input == "" {
		return nil, errors.New("input is empty")
	}

	parts := strings.Fields(input)

	cmd := parts[0]
	cmdLowercase := strings.ToLower(cmd)

	if cmdLowercase == "help" {
		return &Command{
			CmdType: CmdShowHelp,
			Payload: nil,
		}, nil
	}

	if cmdLowercase == "exit" || cmdLowercase == "quit" {
		return &Command{
			CmdType: CmdExit,
			Payload: nil,
		}, nil
	}

	if cmdLowercase == "add" {
		return parseCommandAddPeer(parts[1:])
	}

	if cmdLowercase == "save" {
		return parseCommandSaveFile(parts[1:])
	}

	if cmdLowercase == "ls" || cmdLowercase == "list" {
		return &Command{
			CmdType: CmdList,
			Payload: nil,
		}, nil
	}

	if cmdLowercase == "get" {
		return parseCommandGetFile(parts[1:])
	}

	return nil, fmt.Errorf("unknown command: '%s'", cmd)
}

func parseCommandAddPeer(parts []string) (*Command, error) {
	if len(parts) != 3 {
		return nil, errors.New("expected: <peer name> <peer IP addr> <peer listen port>")
	}

	name := parts[0]
	addrStr := parts[1]
	portStr := parts[2]

	addr := net.ParseIP(addrStr)
	if addr == nil {
		return nil, errors.New("invalid IP address")
	}

	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return nil, err
	}

	payload := PayloadAddPeer{
		Name: name,
		Addr: addr,
		Port: uint16(port),
	}

	return &Command{
		CmdType: CmdAddPeer,
		Payload: payload,
	}, nil
}

func parseCommandSaveFile(parts []string) (*Command, error) {
	if len(parts) != 1 {
		return nil, errors.New("expected: <file path>")
	}

	filepath := parts[0]
	if filepath == "" {
		return nil, errors.New("file path can't be empty")
	}

	paylod := PayloadSaveFile{
		Path: filepath,
	}

	return &Command{
		CmdType: CmdSaveFile,
		Payload: paylod,
	}, nil
}

func parseCommandGetFile(parts []string) (*Command, error) {
	if len(parts) != 1 {
		return nil, errors.New("expected: <file name>")
	}

	filename := parts[0]
	if filename == "" {
		return nil, errors.New("file name can't be empty")
	}

	payload := PayloadGetFile{
		FileName: filename,
	}

	return &Command{
		CmdType: CmdGetFile,
		Payload: payload,
	}, nil
}
