package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"github.com/maximekuhn/diskgo/internal/file"
)

// XXX: thread safety ?

type AESFileEncryptor struct {
	secretKey []byte
	nonces    map[string][]byte
}

func NewAESFileEncryptor(secretKey []byte) *AESFileEncryptor {
	// TODO: validate key
	return &AESFileEncryptor{secretKey: secretKey, nonces: make(map[string][]byte)}
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
	f.Data = ciphertext
	fe.nonces[f.Name] = nonce

	return nil
}

func (fe *AESFileEncryptor) Decrypt(f *file.File) error {
	nonce, ok := fe.nonces[f.Name]
	if !ok {
		return errors.New("no nonce found for file")
	}

	c, err := aes.NewCipher(fe.secretKey)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return err
	}

	plaintext, err := gcm.Open(nil, nonce, f.Data, nil)
	if err != nil {
		return err
	}

	f.Data = plaintext

	return nil
}
