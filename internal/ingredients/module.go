package ingredients

import (
	"database/sql"
	"net/http"
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

func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /ingredients", m.handler.CreateIngredient)
	mux.HandleFunc("GET /ingredients", m.handler.ListIngredients)
	mux.HandleFunc("GET /ingredients/{id}", m.handler.GetIngredient)
	mux.HandleFunc("DELETE /ingredients/", m.handler.DeleteIngredient)
}
