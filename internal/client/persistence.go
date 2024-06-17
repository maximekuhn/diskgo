package client

import (
	"encoding/json"
	"os"
)

// clientPersistence is a structure allowing to store and restore client's state
// on the disk.
type clientPersistence struct {
	clientState   state
	writeFilePath string
}

func newClientPeristence(writeFilePath string) *clientPersistence {
	return &clientPersistence{writeFilePath: writeFilePath}
}

// state is the client's state as it's stored on the disk (JSON format)
type state struct {
	Peers []statePeer `json:"Peers"`
	Files []stateFile `json:"Files"`
}

type statePeer struct {
	Name string `json:"Name"`
	Addr string `json:"Addr"`
}

type stateFile struct {
	Filename string `json:"Filename"`
	Peername string `json:"Peername"`
}

func (m *clientPersistence) writeToDisk() error {
	snapshot, err := json.Marshal(m.clientState)
	if err != nil {
		return err
	}

	outFile, err := os.Create(m.writeFilePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = outFile.Write(snapshot)
	if err != nil {
		return err
	}

	return nil
}

func (m *clientPersistence) readFromDisk() (*state, error) {
	data, err := os.ReadFile(m.writeFilePath)
	if err != nil {
		return nil, err
	}

	var clientState state
	err = json.Unmarshal(data, &clientState)
	if err != nil {
		return nil, err
	}

	return &clientState, nil
}
