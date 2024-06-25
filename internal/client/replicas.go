package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/maximekuhn/diskgo/internal/file"
	"github.com/maximekuhn/diskgo/internal/protocol"
)

type ReplicasStrategy interface {
	// Save a file according to the replicas strategy
	// If the file has been successfully saved, no error is returned
	Save(f *file.File, manager *peersManager, nickname string) error
}

// BasicReplicationManager implements ReplicasStrategy.
// It tries to save the file to replicasCount peers.
// If it's not possible (not enough peer, ...), an error is returned.
type BasicReplicationManager struct {
	replicasCount uint
}

func NewBasicReplicationManager(replicasCount uint) (*BasicReplicationManager, error) {
	if replicasCount <= 0 {
		return nil, errors.New("replicasCount can't be <= 0")
	}
	return &BasicReplicationManager{replicasCount: replicasCount}, nil
}

func (b *BasicReplicationManager) Save(f *file.File, manager *peersManager, nickname string) error {
	peers := manager.getAllKnownPeers()
	if len(peers) < int(b.replicasCount) {
		return errors.New("not enough known peers")
	}

	req := protocol.Message{
		MsgType: protocol.MsgSaveFile,
		From:    nickname,
		Payload: protocol.SaveFileReqPayload{
			File: *f,
		},
	}

	currentCount := 0

	fmt.Printf("%d known peers\n", len(peers))

	for _, p := range peers {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		res, err := sendRequest(ctx, &req, p)
		if err != nil {
			fmt.Printf("peer %s could not save file because '%s'\n", p.Name, err.Error())
			continue
		}

		if res.MsgType != protocol.MsgSaveFileRes {
			continue
		}

		payload, ok := res.Payload.(protocol.SaveFileResPayload)
		if !ok {
			continue
		}

		if !payload.Ok {
			fmt.Printf("peer %s could not save file because '%s'\n", p.Name, payload.Reason)
			continue
		}

		currentCount++
		manager.addFilePeerStorage(f.Name, p)

		fmt.Printf("peer %s saved file\n", p.Name)

		if currentCount == int(b.replicasCount) {
			return nil
		}
	}

	if currentCount != int(b.replicasCount) {
		return fmt.Errorf("could not save to enough peers: %d / %d", currentCount, b.replicasCount)
	}

	return nil
}
