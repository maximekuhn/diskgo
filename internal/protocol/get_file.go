package protocol

import "github.com/maximekuhn/diskgo/internal/file"

type GetFileReqPayload struct {
    FileName string
}

// if Ok is true, then the file is present
type GetFileResPayload struct {
    Ok bool
    File file.File
}
