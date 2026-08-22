package ingredients

import (
	"net/http"
)

type Module struct {
	handler *Handler
}

func NewModule(handler *Handler) *Module {
	return &Module{
		handler: handler,
	}
}

func (mod *Module) RegisterRoutes(
	mux *http.ServeMux,
	auth func(http.Handler) http.Handler,
) {
	mux.Handle("POST /ingredients", auth(http.HandlerFunc(mod.handler.CreateIngredient)))
	mux.Handle("GET /ingredients", auth(http.HandlerFunc(mod.handler.ListIngredients)))
	mux.Handle("GET /ingredients/{id}", auth(http.HandlerFunc(mod.handler.GetIngredient)))
	mux.Handle("DELETE /ingredients/", auth(http.HandlerFunc(mod.handler.DeleteIngredient)))
}
