package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"Feastio/internal/model"
)

// Database Queries
var (
	createIngredient = `
	INSERT INTO ingredients (
		"name", "category", "defaultMeasurementType", "description",
		"createdAt", "createdBy", "modifiedAt", "modifiedBy"
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING *
	`
	getIngredient    = `SELECT * FROM ingredients WHERE id=$1`
	listIngredients  = `SELECT * FROM ingredients`
	deleteIngredient = `DELETE FROM ingredients WHERE id = $1`
)

// Error Definitions
var (
	ErrIngredientNotFound      = errors.New("Ingredient not found")
	ErrIngredientAlreadyExists = errors.New("Ingredient already exists")
)

// Repository Functions
func (r *Repository) CreateIngredient(ctx context.Context, i model.Ingredient) (model.Ingredient, error) {
	var now = time.Now().UTC()
	var ingredient model.Ingredient

	err := r.db.QueryRowContext(
		ctx,
		createIngredient,
		i.Name,
		i.Category,
		i.DefaultMeasurementType,
		i.Description,
		now,
		i.CreatedBy,
		now,
		i.CreatedBy,
	).Scan(
		&ingredient.Id,
		&ingredient.Name,
		&ingredient.Category,
		&ingredient.DefaultMeasurementType,
		&ingredient.Description,
		&ingredient.CreatedAt,
		&ingredient.CreatedBy,
		&ingredient.ModifiedAt,
		&ingredient.ModifiedBy,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Ingredient{}, ErrIngredientNotFound
		}
		return model.Ingredient{}, err
	}

	return ingredient, nil
}

func (r *Repository) GetIngredient(ctx context.Context, id int64) (model.Ingredient, error) {
	var ingredient model.Ingredient

	err := r.db.QueryRowContext(
		ctx,
		getIngredient,
		id,
	).Scan(
		&ingredient.Id,
		&ingredient.Name,
		&ingredient.Category,
		&ingredient.DefaultMeasurementType,
		&ingredient.Description,
		&ingredient.CreatedAt,
		&ingredient.CreatedBy,
		&ingredient.ModifiedAt,
		&ingredient.ModifiedBy,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Ingredient{}, ErrIngredientNotFound
		}
		return model.Ingredient{}, err
	}

	return ingredient, nil
}

func (r *Repository) ListIngredients(ctx context.Context) ([]model.Ingredient, error) {
	rows, err := r.db.QueryContext(
		ctx,
		listIngredients,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []model.Ingredient

	for rows.Next() {
		var i model.Ingredient
		if err := rows.Scan(
			&i.Id,
			&i.Name,
			&i.Category,
			&i.DefaultMeasurementType,
			&i.Description,
			&i.CreatedAt,
			&i.CreatedBy,
			&i.ModifiedAt,
			&i.ModifiedBy,
		); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, i)
	}

	return ingredients, nil
}

func (r *Repository) DeleteIngredient(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(
		ctx,
		deleteIngredient,
		id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrIngredientNotFound
	}

	return nil
}
