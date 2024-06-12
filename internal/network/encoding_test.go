package network

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"testing"

	"github.com/maximekuhn/diskgo/internal/protocol"
)

func TestEncodeHeadersMsgType(t *testing.T) {
	message := protocol.Message{
		MsgType: protocol.MsgGetFile,
		From:    "toto",
		Payload: protocol.GetFileReqPayload{
			FileName: "passwords.txt",
		},
	}

	encodedMsg, err := Encode(message)
	if err != nil {
		t.Fatalf("error while encoding message: %s", err)
	}

	expectedMsgType := message.MsgType
	actualMsgType := protocol.MsgType(encodedMsg[0])
	if expectedMsgType != actualMsgType {
		t.Fatalf("msg type: got %d want %d", actualMsgType, expectedMsgType)
	}
}

func TestEncodeHeadersPayloadLength(t *testing.T) {
	message := protocol.Message{
		MsgType: protocol.MsgGetFile,
		From:    "toto",
		Payload: protocol.GetFileReqPayload{
			FileName: "passwords.txt",
		},
	}

	encodedMsg, err := Encode(message)
	if err != nil {
		t.Fatalf("error while encoding message: %s", err)
	}

	expectedPayloadLength := 28
	actualPayloadLength := binary.BigEndian.Uint32(encodedMsg[5:9])
	if expectedPayloadLength != int(actualPayloadLength) {
		t.Fatalf("payload length: got %d want %d", actualPayloadLength, expectedPayloadLength)
	}
}

func TestEncodeHeadersSendersNicknameLength(t *testing.T) {
	message := protocol.Message{
		MsgType: protocol.MsgGetFile,
		From:    "toto",
		Payload: protocol.GetFileReqPayload{
			FileName: "passwords.txt",
		},
	}

	encodedMsg, err := Encode(message)
	if err != nil {
		t.Fatalf("error while encoding message: %s", err)
	}

	expectedSendersNicknameLength := 4
	actualSendersNicknameLength := binary.BigEndian.Uint32(encodedMsg[1:5])
	if expectedSendersNicknameLength != int(actualSendersNicknameLength) {
		t.Fatalf("senders nickname length: got %d want %d", actualSendersNicknameLength, expectedSendersNicknameLength)
	}
}

func TestEncodePayload(t *testing.T) {
	message := protocol.Message{
		MsgType: protocol.MsgGetFile,
		From:    "toto",
		Payload: protocol.GetFileReqPayload{
			FileName: "passwords.txt",
		},
	}

	encodedMsg, err := Encode(message)
	if err != nil {
		t.Fatalf("error while encoding message: %s", err)
	}

	expectedPayload, err := json.Marshal(message.Payload)
	if err != nil {
		t.Fatalf("error while marshalling payload: %s", err)
	}
	expectedPayload = append([]byte(message.From), expectedPayload...)

	payloadStartIdx := 9
	actualPayload := encodedMsg[payloadStartIdx:]
	if !bytes.Equal(actualPayload, expectedPayload) {
		t.Fatalf("payload: got %v want %v", actualPayload, expectedPayload)
	}
}

func TestEncodeFull(t *testing.T) {
	message := protocol.Message{
		MsgType: protocol.MsgGetFile,
		From:    "toto",
		Payload: protocol.GetFileReqPayload{
			FileName: "passwords.txt",
		},
	}

	encodedMsg, err := Encode(message)
	if err != nil {
		t.Fatalf("error while encoding message: %s", err)
	}

	expectedMsg := []byte{0, 0, 0, 0, 4, 0, 0, 0, 28, 116, 111, 116, 111, 123, 34, 70, 105, 108, 101, 78, 97, 109, 101, 34, 58, 34, 112, 97, 115, 115, 119, 111, 114, 100, 115, 46, 116, 120, 116, 34, 125}

	if !bytes.Equal(encodedMsg, expectedMsg) {
		t.Fatalf("full encoding: got %v want %v", encodedMsg, expectedMsg)
	}
}
