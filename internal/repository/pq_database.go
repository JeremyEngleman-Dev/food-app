package repository

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
	ErrDuplicateUsername = errors.New("Duplicate username found")
	ErrDuplicateEmail    = errors.New("Duplicate email found")
)

var uniqueConstraints = map[string]error{
	"user_unique_name":  ErrDuplicateUsername,
	"user_unique_email": ErrDuplicateEmail,
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
