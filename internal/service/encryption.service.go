package service

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"strings"
)

func EncryptData(data string, key cipher.Block) (string, error) {
	normalized := strings.TrimSpace(strings.ToLower(data))

	aesGCM, err := cipher.NewGCM(key)
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

func DecryptData(encryptedData string, key cipher.Block) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(key)
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
