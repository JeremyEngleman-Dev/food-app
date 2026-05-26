package middleware

import (
	"database/sql"
)

type Middleware struct {
	service *Service
}

func New(db *sql.DB) *Middleware {
	repo := NewRepository(db)
	service := NewService(repo)

	return &Middleware{service: service}
}
