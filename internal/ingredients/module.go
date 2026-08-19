package ingredients

import (
	"database/sql"
	"net/http"

	"foodapp/internal/middleware"
)

type Module struct {
	service *Service
	handler *Handler
}

func New(db *sql.DB) *Module {
	repo := NewRepository(db)
	service := NewService(repo)
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
	mux.Handle("POST /ingredients", mw.Authentication(http.HandlerFunc(m.handler.CreateIngredient)))
	mux.Handle("GET /ingredients", mw.Authentication(http.HandlerFunc(m.handler.ListIngredients)))
	mux.Handle("GET /ingredients/{id}", mw.Authentication(http.HandlerFunc(m.handler.GetIngredient)))
	mux.Handle("DELETE /ingredients/", mw.Authentication(http.HandlerFunc(m.handler.DeleteIngredient)))
}
