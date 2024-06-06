package protocol

type MsgType byte

const (
	MsgGetFile MsgType = iota
    MsgGetFileRes
)

type Message struct {
	MsgType MsgType
	Payload interface{}
}
