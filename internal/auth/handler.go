package auth

import (
	"context"
	"encoding/json"
	"errors"
	"foodapp/internal/logging"
	m "foodapp/internal/models"
	"foodapp/internal/platform/database"
	"foodapp/internal/platform/security"
	"foodapp/internal/platform/utils"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(
	mux *http.ServeMux,
	auth func(http.Handler) http.Handler,
) {
	mux.HandleFunc("POST /auth/login", h.LoginUser)
	mux.Handle("POST /auth/logout", auth(http.HandlerFunc(h.LogoutUser)))
	mux.HandleFunc("POST /auth/refresh", h.RefreshToken)
}

// Handlers
func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer r.Body.Close()

	var login m.LoginInfo
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		utils.HttpJsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	tokens, err := h.service.LoginUser(
		ctx,
		login,
	)
	if err != nil {
		ParseError(w, err)
		return
	}

	var finalWriter http.ResponseWriter = w
	if rw, ok := w.(*logging.ResponseWriter); ok {
		rw.UserCtx = tokens.UserCtx
		finalWriter = rw
	}

	tokenCookie := &http.Cookie{
		Name:     "access_token",
		Value:    tokens.Token,
		Path:     "/",
		HttpOnly: security.HttpOnly,
		Secure:   security.SecureOverHTTPS,
		MaxAge:   int(security.AccessTokenTTL),
	}

	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.Refresh,
		Path:     "/",
		HttpOnly: security.HttpOnly,
		Secure:   security.SecureOverHTTPS,
		MaxAge:   int(security.RefreshTokenTTL),
	}

	http.SetCookie(w, tokenCookie)
	http.SetCookie(w, refreshCookie)
	finalWriter.WriteHeader(http.StatusOK)
}

func (h *Handler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var userCtx m.UserContext
	if rw, ok := w.(*logging.ResponseWriter); !ok {
		utils.HttpJsonResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	} else {
		userCtx = rw.UserCtx
	}

	err := h.service.LogoutUser(ctx, userCtx.UserId)
	if err != nil {
		ParseError(w, err)
		return
	}

	tokenCookie := &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: security.HttpOnly,
		Secure:   security.SecureOverHTTPS,
		MaxAge:   -1,
	}

	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: security.HttpOnly,
		Secure:   security.SecureOverHTTPS,
		MaxAge:   -1,
	}

	http.SetCookie(w, tokenCookie)
	http.SetCookie(w, refreshCookie)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		utils.HttpJsonResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	tokens, err := h.service.RefreshToken(ctx, cookie.Value)
	if err != nil {
		utils.HttpJsonResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	var finalWriter http.ResponseWriter = w

	if rw, ok := w.(*logging.ResponseWriter); ok {
		rw.UserCtx = tokens.UserCtx
		finalWriter = rw
	}

	tokenCookie := &http.Cookie{
		Name:     "access_token",
		Value:    tokens.Token,
		Path:     "/",
		HttpOnly: security.HttpOnly,
		Secure:   security.SecureOverHTTPS,
		MaxAge:   int(security.AccessTokenTTL),
	}

	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.Refresh,
		Path:     "/",
		HttpOnly: security.HttpOnly,
		Secure:   security.SecureOverHTTPS,
		MaxAge:   int(security.RefreshTokenTTL),
	}

	http.SetCookie(w, tokenCookie)
	http.SetCookie(w, refreshCookie)
	finalWriter.WriteHeader(http.StatusOK)
}

// Helper Functions
func ParseError(w http.ResponseWriter, err error) {
	if errors.Is(err, context.Canceled) {
		return
	}

	if errors.Is(err, context.DeadlineExceeded) {
		utils.HttpJsonResponse(w, http.StatusGatewayTimeout, map[string]string{"error": "Gateway timeout"})
		return
	}

	if errors.Is(err, ErrBadCredentials) {
		utils.HttpJsonResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}
	if errors.Is(err, ErrExpiredRefreshToken) {
		utils.HttpJsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Bad Request - Please try logging in again"})
		return
	}

	var appErr *database.AppError
	if errors.As(err, &appErr) {
		switch appErr.Type {
		case database.ErrTypeNotFound:
			utils.HttpJsonResponse(w, http.StatusNotFound, map[string]string{"error": "Invalid email or password"})
			return
		case database.ErrTypeConflict:
			w.WriteHeader(http.StatusOK)
			return
		case database.ErrTypeFailedCreation:
			utils.HttpJsonResponse(w, http.StatusNotFound, map[string]string{"error": "Failed to log in"})
			return
		default:
			utils.HttpJsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			return
		}
	}

	utils.HttpJsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
}
