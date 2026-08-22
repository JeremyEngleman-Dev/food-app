package ingredients

import (
	"context"
	"database/sql"
	"errors"
	m "foodapp/internal/models"
	"foodapp/internal/platform/database"
	"time"
)

type Repository interface {
	CreateIngredient(ctx context.Context, i m.Ingredient) (m.Ingredient, error)
	GetIngredientById(ctx context.Context, id int64) (m.Ingredient, error)
	ListIngredients(ctx context.Context) ([]m.Ingredient, error)
	DeleteIngredient(ctx context.Context, id int64) error
}

type repo struct {
	db *sql.DB
}

func Newrepo(db *sql.DB) Repository {
	return &repo{db: db}
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

// repo Functions
func (r *repo) CreateIngredient(ctx context.Context, i m.Ingredient) (m.Ingredient, error) {
	var now = time.Now().UTC()
	var ingredient m.Ingredient

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
		return m.Ingredient{}, database.DBError(err)
	}

	return ingredient, nil
}

func (r *repo) GetIngredientById(ctx context.Context, id int64) (m.Ingredient, error) {
	var ingredient m.Ingredient

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
		return m.Ingredient{}, database.DBError(err)
	}

	return ingredient, nil
}

func (r *repo) ListIngredients(ctx context.Context) ([]m.Ingredient, error) {
	rows, err := r.db.QueryContext(
		ctx,
		listIngredients,
	)
	if err != nil {
		return nil, database.DBError(err)
	}
	defer rows.Close()

	ingredients := []m.Ingredient{}

	for rows.Next() {
		var i m.Ingredient
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
			return nil, database.DBError(err)
		}
		ingredients = append(ingredients, i)
	}

	return ingredients, nil
}

func (r *repo) DeleteIngredient(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(
		ctx,
		deleteIngredient,
		id,
	)
	if err != nil {
		return database.DBError(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return database.DBError(err)
	}

	if rowsAffected == 0 {
		return &database.AppError{Type: database.ErrTypeNotFound, Err: errors.New("Ingredient not found")}
	}

	return nil
}
