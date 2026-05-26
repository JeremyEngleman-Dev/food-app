package users

import (
	"database/sql"
	"net/http"

	"foodapp/internal/middleware"
	"foodapp/internal/platform/config"
	"foodapp/internal/platform/security"
)

type Module struct {
	service *Service
	handler *Handler
}

func New(db *sql.DB, cfg config.Config) *Module {
	repo := NewRepository(db)
	service := NewService(repo, security.SecurityKeys{
		Encryption: cfg.EmailEncryptionKey,
		Hash:       cfg.EmailHashKey,
	})
	handler := NewHandler(service)

	return &Module{
		service: service,
		handler: handler,
	}
}

func (m *Module) RegisterRoutes(
	mux *http.ServeMux,
	mw *middleware.Middleware,
) {
	mux.Handle("GET /users/{id}", mw.Authentication(mw.Admin(http.HandlerFunc(m.handler.GetUser))))
	mux.HandleFunc("POST /users", m.handler.CreateUser)
	mux.Handle("GET /users", mw.Authentication(mw.Admin(http.HandlerFunc(m.handler.ListUsers))))
	mux.Handle("DELETE /users/{id}", mw.Authentication(mw.Admin(http.HandlerFunc(m.handler.DeleteUser))))
}

func (m *Module) Service() *Service {
	return m.service
}
