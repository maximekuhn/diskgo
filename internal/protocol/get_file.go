package protocol

import "github.com/maximekuhn/diskgo/internal/file"

type GetFileReqPayload struct {
	FileName string
}

// GetFileResPayload
//
// Ok is set to true if the file has been found
type GetFileResPayload struct {
	Ok   bool
	File file.File
}
