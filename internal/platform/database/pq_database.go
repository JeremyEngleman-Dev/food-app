package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/lib/pq"
)

func Open(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	return db, err
}

// Errors
var (
	ErrNotFound = errors.New("Record not found")
	ErrConflict = errors.New("Resource conflict")

	// Users
	ErrDuplicateEmail = errors.New("Duplicate email found")

	// Recipes
	ErrDuplicateRecipe = errors.New("Duplicate recipe found")
)

var uniqueConstraints = map[string]error{
	"user_unique_email":  ErrDuplicateEmail,
	"recipe_unique_name": ErrDuplicateEmail,
}

func DBError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}

	var dbErr *pq.Error
	if errors.As(err, &dbErr) {
		switch dbErr.Code {
		case "23505":
			if constraint, ok := uniqueConstraints[dbErr.Constraint]; ok {
				return constraint
			}
			return ErrConflict
		default:
			return fmt.Errorf("database error: %s: %w", dbErr.Code, dbErr)
		}
	}

	return err
}
