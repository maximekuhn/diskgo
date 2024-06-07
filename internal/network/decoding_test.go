package network

import (
	"bytes"
	"testing"

	"github.com/maximekuhn/diskgo/internal/protocol"
)

func TestDecodeFull(t *testing.T) {
	msg := []byte{0, 0, 0, 0, 28, 123, 34, 70, 105, 108, 101, 78, 97, 109, 101, 34, 58, 34, 112, 97, 115, 115, 119, 111, 114, 100, 115, 46, 116, 120, 116, 34, 125}

	decodedMsg, err := Decode(bytes.NewReader(msg))
	if err != nil {
		t.Fatalf("failed to decode message: %s", err)
	}

	actualMsgType := decodedMsg.MsgType
	expectedMsgType := protocol.MsgGetFile
	if actualMsgType != expectedMsgType {
		t.Fatalf("msg type: got %d want %d", actualMsgType, expectedMsgType)
	}

	payload, ok := decodedMsg.Payload.(protocol.GetFileReqPayload)
	if !ok {
		t.Fatalf("message decoded but payload don't match expectations")
	}

	actualFileName := payload.FileName
	expectedFileName := "passwords.txt"
	if expectedFileName != actualFileName {
		t.Fatalf("filename: got %s want %s", actualFileName, expectedFileName)
	}
}
