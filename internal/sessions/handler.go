package sessions

import (
	m "Feastio/internal/models"
	"Feastio/internal/platform/database"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer r.Body.Close()

	var login m.LoginInfo
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		WriteJsonReturn(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	session, err := h.service.LoginUser(ctx, login)
	if err != nil {
		ParseError(w, err)
		return
	}

	sessionCookie := &http.Cookie{
		Name:     "session_id",
		Value:    session.SessionId,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   3600,
	}

	http.SetCookie(w, sessionCookie)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sessionId, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = h.service.LogoutUser(ctx, sessionId.Value)
	if err != nil {
		ParseError(w, err)
		return
	}

	sessionCookie := &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   -1,
	}

	http.SetCookie(w, sessionCookie)
	w.WriteHeader(http.StatusNoContent)
}

// Error Handling
func WriteJsonReturn(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("failed to write response:", err)
	}
}

func ParseError(w http.ResponseWriter, err error) {
	if errors.Is(err, context.Canceled) {
		return
	}
	if errors.Is(err, context.DeadlineExceeded) {
		WriteJsonReturn(w, http.StatusGatewayTimeout, map[string]string{"error": "Gateway timeout"})
	}

	var appErr *database.AppError
	if errors.As(err, &appErr) {
		switch appErr.Type {
		case database.ErrTypeNotFound:
			fmt.Println("User not found to log in")
			WriteJsonReturn(w, http.StatusNotFound, map[string]string{"error": "Invalid email or password"})
			return
		case database.ErrTypeConflict:
			w.WriteHeader(http.StatusOK)
			return
		case database.ErrTypeFailedCreation:
			fmt.Println("Failed to create session during login")
			WriteJsonReturn(w, http.StatusNotFound, map[string]string{"error": "Failed to log in"})
			return
		default:
			fmt.Println("Unknown database error")
			WriteJsonReturn(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			return
		}
	}

	WriteJsonReturn(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
}
