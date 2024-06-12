package encryption

import "github.com/maximekuhn/diskgo/internal/file"

type FileEncrypter interface {
	// Encrypt file's data
	// If the data could not be encrypted, an error is returned
	Encrypt(f *file.File) error

	// Decrypt file's data
	// If the data  could not be decrypted, an error is returned
	Decrypt(f *file.File) error
}
