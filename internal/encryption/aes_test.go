package encryption

import (
	"bytes"
	"github.com/maximekuhn/diskgo/internal/file"
	"testing"
)

func TestAESFileEncryptor_Encrypt(t *testing.T) {
	data := []byte("gmail: password1245\nwork: superstrongpassword")
	f := file.NewFile("passwords.txt", data)

	secretKey := []byte("i5yrqDhVmvV9YpFBwexikVXYFtC4emd9")
	sut := NewAESFileEncryptor(secretKey)

	err := sut.Encrypt(f)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	if bytes.Equal(data, f.Data) {
		t.Fatal("encrypt failed, file's data not updated")
	}
}

func TestAESFileEncryptor_Decrypt(t *testing.T) {
	data := []byte("gmail: password1245\nwork: superstrongpassword")
	f := file.NewFile("passwords.txt", data)

	secretKey := []byte("i5yrqDhVmvV9YpFBwexikVXYFtC4emd9")
	sut := NewAESFileEncryptor(secretKey)

	err := sut.Encrypt(f)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	// actual test
	err = sut.Decrypt(f)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	if !bytes.Equal(f.Data, data) {
		t.Fatal("decrypt failed, file's data not updated")
	}
}
