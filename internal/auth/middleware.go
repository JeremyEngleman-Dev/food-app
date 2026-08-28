package auth

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	m "foodapp/internal/models"
	"foodapp/internal/platform/security"

	"github.com/golang-jwt/jwt/v5"
)

type Middleware interface {
	Authenticate(next http.Handler) http.Handler
	Admin(next http.Handler) http.Handler
}

type middleware struct {
	token security.Token
}

func NewMiddleware(token security.Token) *middleware {
	return &middleware{
		token: token,
	}
}

func (mw *middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := mw.token.ValidateJWTToken(cookie.Value)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			if errors.Is(err, jwt.ErrTokenExpired) {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		userContext := m.UserContext{
			UserId: claims.UserID,
			Role:   claims.Role,
		}

		ctx := context.WithValue(
			r.Context(),
			"userCtx",
			userContext,
		)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (mw *middleware) Admin(next http.Handler) http.Handler {
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
