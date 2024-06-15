package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"github.com/maximekuhn/diskgo/internal/file"
)

type AESFileEncryptor struct {
	secretKey []byte
}

func NewAESFileEncryptor(secretKey []byte) *AESFileEncryptor {
	// TODO: validate key
	return &AESFileEncryptor{secretKey: secretKey}
}

func (fe *AESFileEncryptor) Encrypt(f *file.File) error {
	c, err := aes.NewCipher(fe.secretKey)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return err
	}

	ciphertext := gcm.Seal(nil, nonce, f.Data, nil)

	// append nonce first
	f.Data = append(nonce, ciphertext...)

	return nil
}

func (fe *AESFileEncryptor) Decrypt(f *file.File) error {
	c, err := aes.NewCipher(fe.secretKey)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return err
	}

	// try to get nonce
	nonceSize := gcm.NonceSize()
	if len(f.Data) < nonceSize {
		return errors.New("ciphertext too short")
	}
	nonce, ciphertext := f.Data[:nonceSize], f.Data[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	f.Data = plaintext

	return nil
}
