package store

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/maximekuhn/diskgo/internal/file"
	"testing"
)

func TestGetFilePath(t *testing.T) {
	tests := []struct {
		filename string
		peername string
		rootDir  string
		want     string
	}{
		{
			"report_final2.zip",
			"bill",
			"~/Documents/peersfiles",
			"~/Documents/peersfiles/e8375d7cd983efcbf956da5937050ffc/fdff96494001fd558da2c2c4f617a082",
		},
		{
			"report_final2.zip",
			"bill",
			"~/Documents/peersfiles/",
			"~/Documents/peersfiles/e8375d7cd983efcbf956da5937050ffc/fdff96494001fd558da2c2c4f617a082",
		},
		{
			"report_final2.zip",
			"john",
			"~/Documents/peersfiles/",
			"~/Documents/peersfiles/527bd5b5d689e2c32ae974c6229ff785/fdff96494001fd558da2c2c4f617a082",
		},
	}

	for _, test := range tests {
		testName := fmt.Sprintf("%s %s %s", test.rootDir, test.filename, test.peername)
		t.Run(testName, func(t *testing.T) {
			got := getPath(test.rootDir, test.filename, test.peername)
			if got != test.want {
				t.Errorf("got %s, want %s", got, test.want)
			}
		})
	}
}

func TestFsFileStore_Save(t *testing.T) {
	rootDir := t.TempDir()

	filename := "shoppings_list.txt"
	data := ""
	f := file.NewFile(filename, []byte(data))

	peername := "toto"

	sut := NewFsFileStore(rootDir, 0)

	err := sut.Save(f, peername)
	if err != nil {
		t.Fatalf("failed to save file: %s", err)
	}
}

func TestFsFileStore_GetNone(t *testing.T) {
	rootDir := t.TempDir()
	sut := NewFsFileStore(rootDir, 0)

	f, err := sut.Get("file.txt", "toto")
	if err == nil {
		if !errors.Is(err, ErrFileNotFound) {
			t.Fatalf("expected ErrFileNotFound, got %s", err)
		}
	}

	if f != nil {
		t.Fatal("should not have been able to get a file")
	}
}

func TestFsFileStore_GetNonePeername(t *testing.T) {
	rootDir := t.TempDir()
	f := file.NewFile("data.csv", []byte("3.14"))
	sut := NewFsFileStore(rootDir, 0)
	err := sut.Save(f, "toto")
	if err != nil {
		t.Fatalf("failed to save file: %s", err)
	}

	// try to get data.csv but as "user" instead of "toto"
	f, err = sut.Get("data.csv", "user")
	if err == nil {
		if !errors.Is(err, ErrFileNotFound) {
			t.Fatalf("expected ErrFileNotFound, got %s", err)
		}
	}

	if f != nil {
		t.Fatal("should not have been able to get a file")
	}
}

func TestFsFileStore_Get(t *testing.T) {
	rootDir := t.TempDir()

	filename := "shoppings_list.txt"
	data := []byte("- ketchup\n- tomatoes")
	f := file.NewFile(filename, data)

	peername := "toto"

	sut := NewFsFileStore(rootDir, 0)

	err := sut.Save(f, peername)
	if err != nil {
		t.Fatalf("failed to save file: %s", err)
	}

	f, err = sut.Get(filename, peername)
	if err != nil {
		t.Fatalf("failed to getfile: %s", err)
	}

	if f == nil {
		t.Fatal("failed to get file (file is nil)")
	}

	if !bytes.Equal(data, f.Data) {
		t.Fatalf("found file but data is not the same: got %s want %s", f.Data, data)
	}
}

func TestFsFileStore_HasEnoughDiskSpace(t *testing.T) {
	tests := []struct {
		maxSizeKB          int64
		currentSizeKB      int64
		fileSizeBytes      int64
		hasEnoughDiskSpace bool
	}{
		{
			maxSizeKB:          1,
			currentSizeKB:      1,
			fileSizeBytes:      10,
			hasEnoughDiskSpace: false,
		},
		{
			maxSizeKB:          1024,
			currentSizeKB:      0,
			fileSizeBytes:      10,
			hasEnoughDiskSpace: true,
		},
		{
			maxSizeKB:          1024,
			currentSizeKB:      999,
			fileSizeBytes:      65536,
			hasEnoughDiskSpace: false,
		},
	}

	rootDir := t.TempDir()

	for _, test := range tests {
		testName := fmt.Sprintf("%d %d %d %v", test.maxSizeKB, test.currentSizeKB, test.fileSizeBytes, test.hasEnoughDiskSpace)
		t.Run(testName, func(t *testing.T) {
			sut := NewFsFileStore(rootDir, test.maxSizeKB)
			sut.currentSizeKB = test.currentSizeKB
			got := sut.hasEnoughDiskSpace(file.NewFile("not important", make([]byte, test.fileSizeBytes)))

			if got != test.hasEnoughDiskSpace {
				t.Errorf("got %v want %v", got, test.hasEnoughDiskSpace)
			}
		})
	}
}

func TestFsFileStore_SaveNoMoreSpace(t *testing.T) {
	rootDir := t.TempDir()
	sut := NewFsFileStore(rootDir, 1)

	// try to save a file that is more than 1 KB
	f := file.NewFile("large_file.txt", make([]byte, 4192))

	err := sut.Save(f, "toto")
	if err == nil {
		t.Fatal("should have failed")
	}

	if !errors.Is(err, ErrNoMoreDiskSpace) {
		t.Fatalf("expected ErrNoMoreDiskSpace, got %s", err)
	}
}

func TestFsFileStore_EnoughDiskSpace(t *testing.T) {
	rootDir := t.TempDir()
	sut := NewFsFileStore(rootDir, 1)

	// try to save a file that is less 1 KB
	f := file.NewFile("large_file.txt", make([]byte, 1023))

	err := sut.Save(f, "toto")
	if err != nil {
		t.Fatalf("failed to save file: %s", err)
	}
}

func TestNewInMemoryFileStore_EnoughThenNotEnoughDiskSpace(t *testing.T) {
	rootDir := t.TempDir()
	sut := NewFsFileStore(rootDir, 1)

	// try to save a file that is less 1 KB
	f := file.NewFile("large_file.txt", make([]byte, 1023))

	err := sut.Save(f, "toto")
	if err != nil {
		t.Fatalf("failed to save file: %s", err)
	}

	// try to save again a new file
	f = file.NewFile("large_file2.txt", make([]byte, 1023))

	err = sut.Save(f, "toto")
	if err == nil {
		t.Fatal("should have failed")
	}
	if !errors.Is(err, ErrNoMoreDiskSpace) {
		t.Fatalf("expected ErrNoMoreDiskSpace, got %s", err)
	}
}
