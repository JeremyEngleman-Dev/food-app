package security

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"strings"
)

type Encryption interface {
	Encrypt(value string) (string, error)
	Decrypt(value string) (string, error)
}

type encryption struct {
	emailEncryptionKey cipher.Block
}

func NewEncryption(emailEncryptionKey cipher.Block) *encryption {
	return &encryption{
		emailEncryptionKey: emailEncryptionKey,
	}
}

func (e *encryption) Encrypt(value string) (string, error) {
	normalized := strings.TrimSpace(strings.ToLower(value))

	aesGCM, err := cipher.NewGCM(e.emailEncryptionKey)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(normalized), nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (e *encryption) Decrypt(value string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(e.emailEncryptionKey)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
