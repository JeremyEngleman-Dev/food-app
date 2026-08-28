package users

import (
	"context"
	"encoding/json"
	"errors"
	mm "foodapp/internal/mappings"
	m "foodapp/internal/models"
	"foodapp/internal/platform/database"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	service     *Service
	userMapping mm.UserMapping
}

func NewHandler(service *Service, userMapping mm.UserMapping) *Handler {
	return &Handler{
		service:     service,
		userMapping: userMapping,
	}
}

func (h *Handler) RegisterRoutes(
	mux *http.ServeMux,
	auth func(http.Handler) http.Handler,
	admin func(http.Handler) http.Handler,
) {
	mux.Handle("GET /users/{id}", auth(http.HandlerFunc(h.GetUser)))
	mux.HandleFunc("POST /users", h.CreateUser)
	mux.Handle("GET /users", auth(http.HandlerFunc(h.ListUsers)))
	mux.Handle("DELETE /users/{id}", auth(admin(http.HandlerFunc(h.DeleteUser))))
}

// Handlers
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer r.Body.Close()

	validRoles := map[string]bool{
		"user":  true,
		"admin": true,
	}

	var request m.CreateUser
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		WriteJsonReturn(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}
	if request.Email == "" || request.Password == "" || !validRoles[request.Role] {
		WriteJsonReturn(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	user, err := h.service.CreateUser(ctx, request)
	if err != nil {
		ParseError(w, err)
		return
	}

	displayUser, err := h.userMapping.ToDisplayUser(user)
	if err != nil {
		WriteJsonReturn(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(displayUser); err != nil {
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

	displayUser, err := h.userMapping.ToDisplayUser(user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(displayUser)
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := h.service.ListUsers(ctx)
	if err != nil {
		ParseError(w, err)
		return
	}

	var displayUsers []m.DisplayUser
	for _, user := range users {
		displayUser, err := h.userMapping.ToDisplayUser(user)
		if err != nil {
			WriteJsonReturn(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			return
		}
		displayUsers = append(displayUsers, displayUser)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(displayUsers)
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
