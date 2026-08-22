package auth

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	m "foodapp/internal/models"
)

type Middleware struct {
	service *Service
}

func NewMiddleware(service *Service) *Middleware {
	return &Middleware{
		service: service,
	}
}

func (mw *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		userContext, err := mw.service.ValidateSession(r.Context(), cookie.Value)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(
			r.Context(),
			"userCtx",
			userContext,
		)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (mw *Middleware) Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userCtx, ok := r.Context().Value("userCtx").(m.UserContext)
		if !ok {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if userCtx.Role != "admin" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
