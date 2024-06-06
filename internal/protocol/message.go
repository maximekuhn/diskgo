package protocol

type MsgType byte

const (
	MsgGetFile MsgType = iota
	MsgGetFileRes

	MsgSaveFile
	MsgSaveFileRes
)

type Message struct {
	MsgType MsgType
	Payload interface{}
}
