package service

import (
	"context"
	"errors"

	"Feastio/internal/dto"
	"Feastio/internal/model"
	"Feastio/internal/repository"
)

// Error Definitions
var (
	ErrIngredientNotFound      = errors.New("Ingredient not found")
	ErrIngredientAlreadyExists = errors.New("Ingredient already exists")
)

// Functions
func (s *Service) CreateIngredient(ctx context.Context, request dto.CreateIngredient) (model.Ingredient, error) {
	i := model.Ingredient{
		Name:                   request.Name,
		Category:               request.Category,
		DefaultMeasurementType: request.DefaultMeasurementType,
		Description:            request.Description,
		CreatedBy:              request.CreatedBy,
	}

	ingredient, err := s.repo.CreateIngredient(ctx, i)
	if err != nil {
		if errors.Is(err, repository.ErrIngredientNotFound) {
			return model.Ingredient{}, ErrIngredientNotFound
		}
		return model.Ingredient{}, err
	}

	return ingredient, nil
}

func (s *Service) GetIngredient(ctx context.Context, id int64) (model.Ingredient, error) {
	ingredient, err := s.repo.GetIngredient(ctx, id)

	if err != nil {
		if errors.Is(err, repository.ErrIngredientNotFound) {
			return model.Ingredient{}, ErrIngredientNotFound
		}
		return model.Ingredient{}, err
	}

	return ingredient, nil
}

func (s *Service) ListIngredients(ctx context.Context) ([]model.Ingredient, error) {
	ingredients, err := s.repo.ListIngredients(ctx)
	if err != nil {
		return nil, err
	}

	return ingredients, nil
}

func (s *Service) DeleteIngredient(ctx context.Context, id int64) error {
	err := s.repo.DeleteIngredient(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrIngredientNotFound) {
			return ErrIngredientNotFound
		}
		return err
	}

	return nil
}
