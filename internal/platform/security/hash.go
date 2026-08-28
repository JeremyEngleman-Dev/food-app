package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Hash interface {
	HashPassword(data string) (string, error)
	ValidateHashPassword(hash, data string) bool
	HashEmail(email string) string
	HashRefreshToken(token []byte) string
}

type hash struct {
	emailHashKey []byte
	tokenHashKey []byte
}

func NewHash(emailHashKey []byte, tokenHashKey []byte) *hash {
	return &hash{
		emailHashKey: emailHashKey,
		tokenHashKey: tokenHashKey,
	}
}

// Passwords

// HashPassword uses bcrypt to hash a plain text password data
func (h *hash) HashPassword(data string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	return string(hash), err
}

// ValidateHashPassword uses bcrypt to validate a plain text password against a hash
func (h *hash) ValidateHashPassword(hash, data string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(data))
	return err == nil
}

// Data
func (h *hash) HashEmail(email string) string {
	hash := hmac.New(sha256.New, h.emailHashKey)
	hash.Write([]byte(strings.ToLower(strings.TrimSpace(email))))
	return base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
}

func (h *hash) HashRefreshToken(token []byte) string {
	sum := sha256.Sum256(token)
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
