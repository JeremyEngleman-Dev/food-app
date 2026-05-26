package sessions

import (
	"database/sql"
	"foodapp/internal/platform/config"
	"foodapp/internal/platform/security"
	"foodapp/internal/users"
	"net/http"
)

type Module struct {
	service *Service
	handler *Handler
}

func New(
	db *sql.DB,
	cfg config.Config,
	users *users.Service,
) *Module {
	repo := NewRepository(db)
	service := NewService(
		repo,
		security.SecurityKeys{
			Encryption: cfg.EmailEncryptionKey,
			Hash:       cfg.EmailHashKey,
		},
		users,
	)
	handler := NewHandler(service)

	return &Module{
		service: service,
		handler: handler,
	}
}

func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /login", m.handler.LoginUser)
	mux.HandleFunc("POST /logout", m.handler.LogoutUser)
}

func (m *Module) Service() *Service {
	return m.service
}
