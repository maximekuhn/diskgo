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
	rootDir       string
	maxSizeKB     int64
	currentSizeKB int64
}

// NewFsFileStore creates a new file system storage with a maximum disk capacity indicated by maxSizeKB. If it's set to 0,
// then no limits are applied.
func NewFsFileStore(rootDir string, maxSizeKB int64) *FsFileStore {
	return &FsFileStore{
		rootDir:       rootDir,
		maxSizeKB:     maxSizeKB,
		currentSizeKB: 0,
	}
}

// Save the given file
func (fs *FsFileStore) Save(f *file.File, peername string) error {
	filepath := getPath(fs.rootDir, f.Name, peername)

	if !fs.hasEnoughDiskSpace(f) {
		return ErrNoMoreDiskSpace
	}

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

func (fs *FsFileStore) hasEnoughDiskSpace(f *file.File) bool {
	if fs.currentSizeKB >= fs.maxSizeKB {
		return false
	}

	fileSizeBytes := int64(len(f.Data))
	currentSizeBytes := fs.currentSizeKB * 1024
	maxSizeBytes := fs.maxSizeKB * 1024

	return currentSizeBytes+fileSizeBytes > maxSizeBytes
}

func getPath(rootDir, filename, peername string) string {
	// get MD5 hash of the filename to avoid needing to sanitaze any filename
	h := md5.New()
	_, err := io.WriteString(h, filename)
	if err != nil {
		panic(err)
	}
	hash := h.Sum(nil)
	hashedFilename := fmt.Sprintf("%x", hash)

	// also use MD5 hash for peer name (same reasons as before)
	h = md5.New()
	_, err = io.WriteString(h, peername)
	if err != nil {
		panic(err)
	}
	hash = h.Sum(nil)
	hashedPeername := fmt.Sprintf("%x", hash)

	return path.Join(rootDir, hashedPeername, hashedFilename)
}
