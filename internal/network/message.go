package network

type MsgType byte

const (
	MsgGetFile MsgType = iota
	MsgGetFileRes

	MsgSaveFile
	MsgSaveFileRes
)

type Message struct {
	Headers Headers
	Payload []byte
}

type Headers struct {
	MsgType       MsgType
	PayloadLength uint32
}
