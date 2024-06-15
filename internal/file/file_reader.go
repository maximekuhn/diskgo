package file

import (
	"os"
	"path/filepath"
)

func ReadFile(path string) (*File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// name of the file is the absolute path
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	name := absolutePath

	f := NewFile(name, data)

	return f, nil
}
