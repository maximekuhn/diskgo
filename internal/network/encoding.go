package network

import (
	"encoding/binary"
	"encoding/json"

	"github.com/maximekuhn/diskgo/internal/protocol"
)

// Encode the given data, so it's ready to be sent over the network
func Encode(msg protocol.Message) ([]byte, error) {
	// create payload
	payload, err := json.Marshal(msg.Payload)
	if err != nil {
		return nil, err
	}

	// create headers (9 bytes)
	// msg type: 1 byte
	// sender nickname's length: 4 bytes
	// payload length: 4 bytes
	senderNicknamesLength := uint32(len(msg.From))
	payloadLength := uint32(len(payload))
	message := make([]byte, 9)
	message[0] = byte(msg.MsgType)
	binary.BigEndian.PutUint32(message[1:5], senderNicknamesLength)
	binary.BigEndian.PutUint32(message[5:], payloadLength)

	// add sender's nickname
	message = append(message, []byte(msg.From)...)

	// add payload
	message = append(message, payload...)

	return message, nil
}
