package users

import (
	"crypto/cipher"
	"database/sql"
	"net/http"

	"Feastio/internal/platform/config"
)

type Module struct {
	service *Service
	handler *Handler
}

type UserKeys struct {
	encryption cipher.Block
	hash       []byte
}

func New(db *sql.DB, cfg config.Config) *Module {
	repo := NewRepository(db)
	service := NewService(repo, UserKeys{
		encryption: cfg.EmailEncryptionKey,
		hash:       cfg.EmailHashKey,
	})
	handler := NewHandler(service)

	return &Module{
		service: service,
		handler: handler,
	}
}

func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /users/{id}", m.handler.GetUser)
	mux.HandleFunc("POST /users", m.handler.CreateUser)
	mux.HandleFunc("GET /users", m.handler.ListUsers)
	mux.HandleFunc("DELETE /users/{id}", m.handler.DeleteUser)
}
