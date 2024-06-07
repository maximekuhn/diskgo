package file

import "os"

func ReadFile(path string) (*File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// for now, name is the path
	name := path
	f := NewFile(name, data)

	return f, nil
}
