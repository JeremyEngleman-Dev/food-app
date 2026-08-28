package security

import (
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	m "foodapp/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

type Token interface {
	CreateJwtToken(user m.User) (string, error)
	CreateRefreshToken(n int) ([]byte, error)
	ValidateJWTToken(tokenString string) (*m.Claims, error)
}

type token struct {
	TokenHashKey []byte
}

func NewToken(TokenHashKey []byte) *token {
	return &token{
		TokenHashKey: TokenHashKey,
	}
}

var (
	ErrTokenExpired = errors.New("JWT token has expired")
)

func (t *token) CreateJwtToken(user m.User) (string, error) {
	claims := m.Claims{
		UserID: user.Id,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(t.TokenHashKey)
}

func (t *token) CreateRefreshToken(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

func (t *token) ValidateJWTToken(tokenString string) (*m.Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&m.Claims{},
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method")
			}

			return t.TokenHashKey, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*m.Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
