package handlers

import (
	"errors"
	"io"

	"github.com/maximekuhn/diskgo/internal/protocol"
	"github.com/maximekuhn/diskgo/internal/store"
)

func HandleSaveFile(msg *protocol.Message, fStore store.FileStore, w io.Writer) error {
	if msg.MsgType != protocol.MsgSaveFile {
		return errors.New("bad message type")
	}
	req, ok := msg.Payload.(protocol.SaveFileReqPayload)
	if !ok {
		return errors.New("correct message type but unepexted payload type")
	}

	err := fStore.Save(&req.File)
	if err != nil {
		return writeResponse(protocol.MsgSaveFileRes, protocol.SaveFileResPayload{
			Ok: false,
		}, w)
	}

	return writeResponse(protocol.MsgSaveFileRes, protocol.SaveFileResPayload{
		Ok: true,
	}, w)
}
