package users

import (
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Password
func HashPassword(data string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	return string(hash), err
}

func CheckHashPassword(hash, data string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(data))
	return err == nil
}

// Encryption
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

// Predictive Hashing
func HashData(data string, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(data))

	return hex.EncodeToString(mac.Sum(nil))
}
