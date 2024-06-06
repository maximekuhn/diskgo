package store

import (
	"errors"

	"github.com/maximekuhn/diskgo/internal/file"
)

var (
	ErrFileNotFound = errors.New("file not found")
)

type FileStore interface {
	// save the given file
	Save(*file.File) error

	// get the given file by name
	// If the file is not found, an error is returned
	Get(string) (*file.File, error)
}
