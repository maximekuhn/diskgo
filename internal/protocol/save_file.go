package protocol

import "github.com/maximekuhn/diskgo/internal/file"

type SaveFileReqPayload struct {
	File file.File
}

type SaveFileResPayload struct {
	Ok bool
}
