package handlers

import (
	"io"

	"github.com/maximekuhn/diskgo/internal/network"
	"github.com/maximekuhn/diskgo/internal/protocol"
)

func writeResponse(msgType protocol.MsgType, payload interface{}, w io.Writer) error {
	encodedMsg, err := network.Encode(msgType, payload)
	if err != nil {
		return err
	}

	// TODO: handle n write
	_, err = w.Write(encodedMsg)
	if err != nil {
		return err
	}

	return nil
}
