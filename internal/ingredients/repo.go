package ingredients

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

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

// Repository Functions
func (r *Repository) CreateIngredient(ctx context.Context, i Ingredient) (Ingredient, error) {
	var now = time.Now().UTC()
	var ingredient Ingredient

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
			return Ingredient{}, ErrIngredientNotFound
		}
		return Ingredient{}, err
	}

	return ingredient, nil
}

func (r *Repository) GetIngredient(ctx context.Context, id int64) (Ingredient, error) {
	var ingredient Ingredient

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
			return Ingredient{}, ErrIngredientNotFound
		}
		return Ingredient{}, err
	}

	return ingredient, nil
}

func (r *Repository) ListIngredients(ctx context.Context) ([]Ingredient, error) {
	rows, err := r.db.QueryContext(
		ctx,
		listIngredients,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []Ingredient

	for rows.Next() {
		var i Ingredient
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
