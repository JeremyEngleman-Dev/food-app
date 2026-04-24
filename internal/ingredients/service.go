package ingredients

import (
	"Feastio/internal/platform/database"
	"context"
	"errors"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Error Definitions
var (
	ErrIngredientNotFound      = errors.New("Ingredient not found")
	ErrIngredientAlreadyExists = errors.New("Ingredient already exists")
)

// Functions
func (s *Service) CreateIngredient(ctx context.Context, request CreateIngredient) (Ingredient, error) {
	i := Ingredient{
		Name:                   request.Name,
		Category:               request.Category,
		DefaultMeasurementType: request.DefaultMeasurementType,
		Description:            request.Description,
		CreatedBy:              request.CreatedBy,
	}

	ingredient, err := s.repo.CreateIngredient(ctx, i)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return Ingredient{}, ErrIngredientNotFound
		}
		return Ingredient{}, err
	}

	return ingredient, nil
}

func (s *Service) GetIngredient(ctx context.Context, id int64) (Ingredient, error) {
	ingredient, err := s.repo.GetIngredient(ctx, id)

	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return Ingredient{}, ErrIngredientNotFound
		}
		return Ingredient{}, err
	}

	return ingredient, nil
}

func (s *Service) ListIngredients(ctx context.Context) ([]Ingredient, error) {
	ingredients, err := s.repo.ListIngredients(ctx)
	if err != nil {
		return nil, err
	}

	return ingredients, nil
}

func (s *Service) DeleteIngredient(ctx context.Context, id int64) error {
	err := s.repo.DeleteIngredient(ctx, id)
	if err != nil {
		if errors.Is(err, ErrIngredientNotFound) {
			return ErrIngredientNotFound
		}
		return err
	}

	return nil
}
