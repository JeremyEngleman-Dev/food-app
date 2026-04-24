package users

import "time"

type User struct {
	Id             int64
	DisplayName    string
	PasswordHash   string
	EmailHash      string
	EmailEncrypted string
	CreatedAt      time.Time
	ModifiedAt     time.Time
	Role           string
}

// Requests
type CreateUser struct {
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	Role        string `json:"role"`
}

// Responses
type DisplayUser struct {
	Id          int64     `json:"id"`
	DisplayName string    `json:"displayName"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"createdAt"`
	ModifiedAt  time.Time `json:"modifiedAt"`
	Role        string    `json:"role"`
}
