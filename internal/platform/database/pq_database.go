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
type ErrorType int

const (
	ErrTypeUnknown ErrorType = iota
	ErrTypeNotFound
	ErrTypeConflict
	ErrTypeDatabase
	ErrTypeFailedCreation
)

type AppError struct {
	Type ErrorType
	Err  error
}

func (e *AppError) Error() string { return e.Err.Error() }
func (e *AppError) Unwrap() error { return e.Err }

func DBError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return &AppError{Type: ErrTypeNotFound, Err: err}
	}

	var dbErr *pq.Error
	if errors.As(err, &dbErr) {
		switch dbErr.Code {
		case "23505":
			return &AppError{Type: ErrTypeConflict, Err: err}
		default:
			return &AppError{Type: ErrTypeDatabase, Err: fmt.Errorf("database error: %s: %w", dbErr.Code, dbErr)}
		}
	}

	return &AppError{Type: ErrTypeUnknown, Err: err}
}
