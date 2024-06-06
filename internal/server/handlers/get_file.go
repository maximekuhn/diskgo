package handlers

import (
	"errors"

	"github.com/maximekuhn/diskgo/internal/file"
	"github.com/maximekuhn/diskgo/internal/protocol"
	"github.com/maximekuhn/diskgo/internal/store"
)

func HandleGetFile(msg *protocol.Message, fStore store.FileStore) (*protocol.GetFileResPayload, error) {
	if msg.MsgType != protocol.MsgGetFile {
		return nil, errors.New("incorrect message type")
	}

	req, ok := msg.Payload.(protocol.GetFileReqPayload)
	if !ok {
		return nil, errors.New("correct message type but incorrect payload")
	}

	f, err := fStore.Get(req.FileName)
	if err != nil {
		if errors.Is(err, store.ErrFileNotFound) {
			return &protocol.GetFileResPayload{
				Ok:   false,
				File: file.File{},
			}, nil
		}
		return nil, err
	}

	return &protocol.GetFileResPayload{
		Ok:   true,
		File: *f,
	}, nil
}
