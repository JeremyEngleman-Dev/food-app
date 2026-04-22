package model

import "time"

type User struct {
	Id             int64
	Username       string
	PasswordHash   string
	EmailEncrypted string
	CreatedAt      time.Time
	ModifiedAt     time.Time
	Role           string
	DisplayName    string
}
