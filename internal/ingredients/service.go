package ingredients

import (
	"context"
	"errors"
	m "foodapp/internal/models"
	"foodapp/internal/platform/database"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Functions
func (s *Service) CreateIngredient(ctx context.Context, request m.CreateIngredient) (m.Ingredient, error) {
	i := m.Ingredient{
		Name:                   request.Name,
		Category:               request.Category,
		DefaultMeasurementType: request.DefaultMeasurementType,
		Description:            request.Description,
		CreatedBy:              request.CreatedBy,
	}

	ingredient, err := s.repo.CreateIngredient(ctx, i)
	if err != nil {
		var appErr *database.AppError
		if errors.As(err, &appErr) {
			if appErr.Type == database.ErrTypeNotFound {
				return m.Ingredient{}, &database.AppError{Type: database.ErrTypeFailedCreation, Err: err}
			}
		}
		return m.Ingredient{}, err
	}

	return ingredient, nil
}

func (s *Service) GetIngredient(ctx context.Context, id int64) (m.Ingredient, error) {
	ingredient, err := s.repo.GetIngredient(ctx, id)

	if err != nil {
		return m.Ingredient{}, err
	}

	return ingredient, nil
}

func (s *Service) ListIngredients(ctx context.Context) ([]m.Ingredient, error) {
	ingredients, err := s.repo.ListIngredients(ctx)
	if err != nil {
		return nil, err
	}

	return ingredients, nil
}

func (s *Service) DeleteIngredient(ctx context.Context, id int64) error {
	err := s.repo.DeleteIngredient(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
