package ingredients

import (
	"context"
	"encoding/json"
	"errors"
	"foodapp/internal/logging"
	m "foodapp/internal/models"
	"foodapp/internal/platform/database"
	"foodapp/internal/platform/utils"
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

func (h *Handler) RegisterRoutes(
	mux *http.ServeMux,
	auth func(http.Handler) http.Handler,
) {
	mux.Handle("POST /ingredients", auth(http.HandlerFunc(h.CreateIngredient)))
	mux.Handle("GET /ingredients", auth(http.HandlerFunc(h.ListIngredients)))
	mux.Handle("GET /ingredients/{id}", auth(http.HandlerFunc(h.GetIngredient)))
	mux.Handle("DELETE /ingredients/", auth(http.HandlerFunc(h.DeleteIngredient)))
}

// Handlers
func (h *Handler) CreateIngredient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var userCtx m.UserContext
	if rw, ok := w.(*logging.ResponseWriter); !ok {
		utils.HttpJsonResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	} else {
		userCtx = rw.UserCtx
	}

	defer r.Body.Close()

	var request m.CreateIngredient
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.HttpJsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}
	if request.Name == "" || request.Description == "" {
		utils.HttpJsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	request.CreatedBy = userCtx.UserId

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
		utils.HttpJsonResponse(w, http.StatusNotFound, map[string]string{"error": "Invalid id"})
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
		utils.HttpJsonResponse(w, http.StatusNotFound, map[string]string{"error": "Invalid id"})
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
func ParseError(w http.ResponseWriter, err error) {
	if errors.Is(err, context.Canceled) {
		return
	}
	if errors.Is(err, context.DeadlineExceeded) {
		utils.HttpJsonResponse(w, http.StatusGatewayTimeout, map[string]string{"error": "Gateway timeout"})
	}

	var appErr *database.AppError
	if errors.As(err, &appErr) {
		switch appErr.Type {
		case database.ErrTypeNotFound:
			utils.HttpJsonResponse(w, http.StatusNotFound, map[string]string{"error": "Ingredient not found"})
			return
		case database.ErrTypeConflict:
			utils.HttpJsonResponse(w, http.StatusConflict, map[string]string{"error": "Ingredient already exist"})
			return
		case database.ErrTypeFailedCreation:
			utils.HttpJsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create Ingredient"})
			return
		default:
			utils.HttpJsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			return
		}
	}

	utils.HttpJsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
}
