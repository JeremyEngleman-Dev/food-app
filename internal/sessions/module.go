package sessions

import (
	"net/http"
)

type Module struct {
	handler *Handler
}

func NewModule(
	handler *Handler,
) *Module {
	return &Module{
		handler: handler,
	}
}

func (m *Module) RegisterRoutes(
	mux *http.ServeMux,
) {
	mux.HandleFunc("POST /login", m.handler.LoginUser)
	mux.HandleFunc("POST /logout", m.handler.LogoutUser)
}
