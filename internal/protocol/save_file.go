package protocol

import "github.com/maximekuhn/diskgo/internal/file"

type SaveFileReqPayload struct {
	File file.File
}

// SaveFileResPayload
//
// # If the file has been saved, Ok is set to true
//
// # If the file has not been saved, Ok is set to false and Reason is filled
//
// Note: Reason can be an empty string if it's an internal error that should not be provided to the client
type SaveFileResPayload struct {
	Ok     bool
	Reason string
}
