package network

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"

	"github.com/maximekuhn/diskgo/internal/protocol"
)

// decode the incoming data (from the network) into a protocol message
func Decode(r io.Reader) (*protocol.Message, error) {
	// read headers (1 byte for msg type, 4 bytes for payload length)
	headersBuf := make([]byte, 5)
	n, err := r.Read(headersBuf)
	if err != nil {
		return nil, err
	}
	if n != 5 {
		return nil, errors.New("did not read enough bytes")
	}

	// check msg type
	msgType := headersBuf[0]
	if msgType > 1 {
		return nil, errors.New("unknown message type")
	}

	// read payload
	payloadLength := binary.BigEndian.Uint32(headersBuf[1:])
	payloadBuf := make([]byte, payloadLength)
	n, err = r.Read(payloadBuf)
	if err != nil {
		return nil, err
	}
	if n != int(payloadLength) {
		return nil, errors.New("did not read enough bytes")
	}

	// check the payload type
	protocolMsgType := protocol.MsgType(msgType)
	payload := payloadBuf[:n]

	if protocolMsgType == protocol.MsgGetFile {
		var pload protocol.GetFileReqPayload
		decoder := json.NewDecoder(bytes.NewReader(payload))
		err = decoder.Decode(&pload)
		if err != nil {
			return nil, err
		}

		return &protocol.Message{
			MsgType: protocolMsgType,
			Payload: pload,
		}, nil
	}

    if protocolMsgType == protocol.MsgGetFileRes {
		var pload protocol.GetFileResPayload
		decoder := json.NewDecoder(bytes.NewReader(payload))
		err = decoder.Decode(&pload)
		if err != nil {
			return nil, err
		}

		return &protocol.Message{
			MsgType: protocolMsgType,
			Payload: pload,
		}, nil
    }

	return nil, errors.New("unknown payload / not implemented yet")
}
