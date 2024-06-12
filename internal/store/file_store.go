package store

import (
	"errors"

	"github.com/maximekuhn/diskgo/internal/file"
)

var (
	ErrFileNotFound = errors.New("file not found")
)

type FileStore interface {
	// Save the given file
	Save(f *file.File, peername string) error

	// Get the given file by name
	// If the file is not found, an error is returned
	Get(filename string, peername string) (*file.File, error)
}
