package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

// AES-256-GCM with 32-byte key
func Encrypt(plain string) (string, error) {
	key := []byte(os.Getenv("JWT_SECRET"))
	if len(key) != 32 {
		return "", errors.New("JWT_SECRET must be 32 bytes for AES-256")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	cipher := gcm.Seal(nonce, nonce, []byte(plain), nil)
	return base64.StdEncoding.EncodeToString(cipher), nil
}

func Decrypt(cipherText string) (string, error) {
	key := []byte(os.Getenv("JWT_SECRET"))
	if len(key) != 32 {
		return "", errors.New("JWT_SECRET must be 32 bytes for AES-256")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", errors.New("invalid cipher")
	}
	nonce, cipher := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plain, err := gcm.Open(nil, nonce, cipher, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}