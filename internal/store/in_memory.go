package store

import (
	"sync"

	"github.com/maximekuhn/diskgo/internal/file"
)

// A dummy implementation of FileStore that saves all files in memory
// This implementation is thread safe
type InMemoryFileStore struct {
	mu    sync.Mutex
	files map[string]*file.File
}

func NewInMemoryFileStore() *InMemoryFileStore {
	return &InMemoryFileStore{
		mu:    sync.Mutex{},
		files: make(map[string]*file.File),
	}
}

// save the given file
func (in *InMemoryFileStore) Save(f *file.File) error {
	in.mu.Lock()
	defer in.mu.Unlock()

	in.files[f.Name] = f

	return nil
}

// get the given file by name
// If the file is not found, an error is returned
func (in *InMemoryFileStore) Get(filename string) (*file.File, error) {
	in.mu.Lock()
	defer in.mu.Unlock()

	f, ok := in.files[filename]
	if !ok {
		return nil, ErrFileNotFound
	}
	if f == nil {
		return nil, ErrFileNotFound
	}

	return f, nil
}
