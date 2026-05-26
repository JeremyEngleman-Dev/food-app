package ingredients

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

// Handlers
func (h *Handler) CreateIngredient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer r.Body.Close()

	var request m.CreateIngredient
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		WriteJsonReturn(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	ingredient, err := h.service.CreateIngredient(ctx, request)
	if err != nil {
		ParseError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ingredient)
}

func (h *Handler) GetIngredient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := strings.TrimPrefix(r.URL.Path, "/ingredients/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteJsonReturn(w, http.StatusNotFound, map[string]string{"error": "Invalid id"})
		return
	}

	ingredient, err := h.service.GetIngredient(ctx, id)
	if err != nil {
		ParseError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ingredient)
}

func (h *Handler) ListIngredients(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ingredients, err := h.service.ListIngredients(ctx)
	if err != nil {
		ParseError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ingredients)
}

func (h *Handler) DeleteIngredient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := strings.TrimPrefix(r.URL.Path, "/ingredients/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteJsonReturn(w, http.StatusNotFound, map[string]string{"error": "Invalid id"})
		return
	}

	err = h.service.DeleteIngredient(ctx, id)
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
			WriteJsonReturn(w, http.StatusNotFound, map[string]string{"error": "Ingredient not found"})
			return
		case database.ErrTypeConflict:
			WriteJsonReturn(w, http.StatusConflict, map[string]string{"error": "Ingredient already exist"})
			return
		case database.ErrTypeFailedCreation:
			WriteJsonReturn(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create Ingredient"})
			return
		default:
			WriteJsonReturn(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			return
		}
	}

	WriteJsonReturn(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
}
