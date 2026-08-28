package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RefreshToken struct {
	RefreshTokenHash string    `json:"refreshTokenHash"`
	UserId           int64     `json:"userId"`
	CreatedAt        time.Time `json:"createdAt"`
	ExpiresAt        time.Time `json:"expiresAt"`
	RevokedAt        time.Time `json:"revokedAt"`
}

type Claims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type ResponseTokens struct {
	Token   string
	Refresh string
}

type UpdateRefreshToken struct {
	RefreshTokenHash *string    `json:"refreshTokenHash"`
	UserId           *int64     `json:"userId"`
	CreatedAt        *time.Time `json:"createdAt"`
	ExpiresAt        *time.Time `json:"expiresAt"`
	RevokedAt        *time.Time `json:"revokedAt"`
}
