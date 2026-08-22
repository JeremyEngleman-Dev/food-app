package models

import "time"

type Session struct {
	SessionId string    `json:"sessionId"`
	UserId    int64     `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type UpdateSessionRequest struct {
	SessionId *string    `json:"sessionId"`
	UserId    *int64     `json:"userId"`
	CreatedAt *time.Time `json:"createdAt"`
	ExpiresAt *time.Time `json:"expiresAt"`
}
