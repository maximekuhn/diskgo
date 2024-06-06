package network

import (
	"encoding/binary"
	"encoding/json"

	"github.com/maximekuhn/diskgo/internal/protocol"
)

// encode the given data so it's ready to be sent over the network
func Encode(msgType protocol.MsgType, data interface{}) ([]byte, error) {
	// create payload
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// create headers (msg type + payload length)
	payloadLength := uint32(len(payload))
	nMsgType := MsgType(msgType)
	message := make([]byte, 5)
	message[0] = byte(nMsgType)
	binary.BigEndian.PutUint32(message[1:], payloadLength)

	// add payload
	message = append(message, payload...)

	return message, nil
}
