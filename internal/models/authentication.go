package models

import "time"

type Session struct {
	SessionId string    `json:"sessionId"`
	UserId    int64     `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
	Role      string    `json:"access"`
}

type LoginInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserContext struct {
	UserId int64  `json:"userId"`
	Role   string `json:"role"`
}
