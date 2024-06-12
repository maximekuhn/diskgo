package store

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/maximekuhn/diskgo/internal/file"
)

// FsFileStore File store that uses the filesystem
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

// Save the given file
func (fs *FsFileStore) Save(f *file.File, peername string) error {
	filepath := getPath(fs.rootDir, f.Name, peername)

	// create file directory if it doesn't exist yet
	filedir := path.Dir(filepath)
	if err := os.MkdirAll(filedir, 0750); err != nil {
		return err
	}

	fmt.Println("filepath", filepath)

	outFile, err := os.Create(filepath)
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

// Get the given file by name
// If the file is not found, an error is returned
func (fs *FsFileStore) Get(filename string, peername string) (*file.File, error) {
	filepath := getPath(fs.rootDir, filename, peername)
	f, err := file.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	// change file name to what the client requested (instead of the MD5 hash)
	f.Name = filename

	return f, nil
}

func getPath(rootDir, filename, peername string) string {
	// get MD5 hash of the filename to avoid needing to sanitaze any filename
	h := md5.New()
	io.WriteString(h, filename)
	hash := h.Sum(nil)
	hashedFilename := fmt.Sprintf("%x", hash)

	// also use MD5 hash for peer name (same reasons as before)
	h = md5.New()
	io.WriteString(h, peername)
	hash = h.Sum(nil)
	hashedPeername := fmt.Sprintf("%x", hash)

	return path.Join(rootDir, hashedPeername, hashedFilename)
}
