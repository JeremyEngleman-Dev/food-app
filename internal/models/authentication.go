package models

type LoginInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserContext struct {
	UserId int64  `json:"userId"`
	Role   string `json:"role"`
}
