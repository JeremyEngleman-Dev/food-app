package users

import (
	"net/http"
)

type Module struct {
	handler *Handler
}

func NewModule(handler *Handler) *Module {
	return &Module{handler: handler}
}

func (mod *Module) RegisterRoutes(
	mux *http.ServeMux,
	auth func(http.Handler) http.Handler,
	admin func(http.Handler) http.Handler,
) {
	mux.Handle("GET /users/{id}", auth(http.HandlerFunc(mod.handler.GetUser)))
	mux.HandleFunc("POST /users", mod.handler.CreateUser)
	mux.Handle("GET /users", auth(http.HandlerFunc(mod.handler.ListUsers)))
	mux.Handle("DELETE /users/{id}", auth(admin(http.HandlerFunc(mod.handler.DeleteUser))))
}
