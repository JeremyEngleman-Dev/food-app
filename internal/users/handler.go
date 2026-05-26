package users

import (
	"context"
	"encoding/json"
	"errors"
	m "foodapp/internal/models"
	"foodapp/internal/platform/database"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer r.Body.Close()

	var request m.CreateUser
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		WriteJsonReturn(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	user, err := h.service.CreateUser(ctx, request)
	if err != nil {
		ParseError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Println("failed to write response:", err)
	}
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteJsonReturn(w, http.StatusNotFound, map[string]string{"error": "Invalid id"})
		return
	}

	user, err := h.service.GetUserById(ctx, id)
	if err != nil {
		ParseError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := h.service.ListUsers(ctx)
	if err != nil {
		ParseError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteJsonReturn(w, http.StatusNotFound, map[string]string{"error": "Invalid id"})
		return
	}

	err = h.service.DeleteUserById(ctx, id)
	if err != nil {
		ParseError(w, err)
		return
	}

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
			WriteJsonReturn(w, http.StatusNotFound, map[string]string{"error": "User not found"})
			return
		case database.ErrTypeConflict:
			WriteJsonReturn(w, http.StatusConflict, map[string]string{"error": "User already exist"})
			return
		case database.ErrTypeFailedCreation:
			WriteJsonReturn(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
			return
		default:
			WriteJsonReturn(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			return
		}
	}

	WriteJsonReturn(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
}
