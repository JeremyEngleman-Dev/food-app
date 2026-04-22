package dto

import (
	"time"
)

// Requests
type CreateUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// Responses
type DisplayUser struct {
	Id         int64     `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
	Role       string    `json:"role"`
}
