package store

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/maximekuhn/diskgo/internal/file"
)

// File store that uses the filesystem
//
// All files from remote peers are stored within the provided rootDir
type FsFileStore struct {
	rootDir string
}

func NewFsFileStore(rootDir string) *FsFileStore {
	return &FsFileStore{
		rootDir: rootDir,
	}
}

// save the given file
func (fs *FsFileStore) Save(f *file.File) error {
	// TODO: create rootDir and maybe sub dirs, if any

	path := getPath(fs.rootDir, f.Name)

	fmt.Println("path", path)

	outFile, err := os.Create(path)
	if err != nil {
		fmt.Println("err create", err)
		return err
	}
	defer outFile.Close()

	_, err = outFile.Write(f.Data)
	if err != nil {
		fmt.Println("err write", err)
		return err
	}

	return nil
}

// get the given file by name
// If the file is not found, an error is returned
func (fs *FsFileStore) Get(filename string) (*file.File, error) {
	path := getPath(fs.rootDir, filename)
	f, err := file.ReadFile(path)
	if err != nil {
		return nil, err
	}

    // change file name to what the client requested (instead of the MD5 hash)
	f.Name = filename

	return f, nil
}

func getPath(rootDir, filename string) string {
	h := md5.New()
	io.WriteString(h, filename)
	hash := h.Sum(nil)
	hashedFilename := fmt.Sprintf("%x", hash)

	return path.Join(rootDir, hashedFilename)
}
