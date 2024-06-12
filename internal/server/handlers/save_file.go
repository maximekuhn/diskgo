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

	err := fStore.Save(&req.File, msg.From)
	if err != nil {
		res := protocol.Message{
			MsgType: protocol.MsgSaveFileRes,
			From:    "todo",
			Payload: protocol.SaveFileResPayload{Ok: false},
		}
		return writeResponse(res, w)
	}

	res := protocol.Message{
		MsgType: protocol.MsgSaveFileRes,
		From:    "todo",
		Payload: protocol.SaveFileResPayload{Ok: true},
	}
	return writeResponse(res, w)
}
