package handlers

import (
	"errors"
	"io"

	"github.com/maximekuhn/diskgo/internal/file"
	"github.com/maximekuhn/diskgo/internal/protocol"
	"github.com/maximekuhn/diskgo/internal/store"
)

func HandleGetFile(msg *protocol.Message, fStore store.FileStore, w io.Writer) error {
	if msg.MsgType != protocol.MsgGetFile {
		return errors.New("incorrect message type")
	}

	req, ok := msg.Payload.(protocol.GetFileReqPayload)
	if !ok {
		return errors.New("correct message type but incorrect payload")
	}

	f, err := fStore.Get(req.FileName, msg.From)
	if err != nil {
		if errors.Is(err, store.ErrFileNotFound) {
			res := protocol.Message{
				MsgType: protocol.MsgGetFileRes,
				From:    "todo",
				Payload: protocol.GetFileResPayload{
					Ok:   false,
					File: file.File{},
				},
			}
			return writeResponse(res, w)
		}

		return err
	}

	res := protocol.Message{
		MsgType: protocol.MsgGetFileRes,
		From:    "todo",
		Payload: protocol.GetFileResPayload{
			Ok:   true,
			File: *f,
		},
	}
	return writeResponse(res, w)
}
