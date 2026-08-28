package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

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

func PatchQueryBuilder(table string, updates map[string]any, identifier string, id any) (string, []any) {
	sets := make([]string, 0, len(updates))
	args := make([]any, 0, len(updates))

	i := 1

	for param, value := range updates {
		sets = append(sets, fmt.Sprintf(`"%s" = $%d`, param, i))
		args = append(args, value)
		i++
	}

	query := fmt.Sprintf(
		`UPDATE %s SET %s WHERE "%s" = $%d`,
		table,
		strings.Join(sets, ", "),
		identifier,
		i,
	)

	args = append(args, id)

	return query, args
}

func PatchMapUpdates(req any) map[string]any {
	updates := make(map[string]any)

	v := reflect.ValueOf(req)
	t := reflect.TypeOf(req)

	if v.Kind() == reflect.Pointer {
		v = v.Elem()
		t = t.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if field.Kind() != reflect.Pointer || field.IsNil() {
			continue
		}

		name := fieldType.Tag.Get("json")

		if name == "" {
			name = fieldType.Name
		}

		// Remove ",omitempty" if present
		if comma := strings.Index(name, ","); comma != -1 {
			name = name[:comma]
		}

		updates[name] = field.Elem().Interface()
	}

	return updates
}
