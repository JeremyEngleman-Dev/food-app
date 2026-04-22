package repository

import (
	"context"
	"errors"
	"time"

	"Feastio/internal/model"
)

// Database Queries
var (
	createUser = `
	INSERT INTO users (
		"username", "passwordHash", "emailEncrypted",
		"createdAt", "modifiedAt",
		"role"
	) VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING "id", "username", "emailEncrypted", "createdAt", "modifiedAt", "role"
	`
	getUser = `
	SELECT "id", "username", "emailEncrypted", "createdAt", "modifiedAt","role"
	FROM users WHERE id=$1
	`
	listUsers = `SELECT "id", "username", "emailEncrypted", "createdAt", "modifiedAt","role"
	FROM users`
	deleteUser = `DELETE FROM users WHERE id = $1`
)

// Error Definitions
var (
	ErrUserNotFound = errors.New("User not found")
)

// Functions
func (r *Repository) CreateUser(ctx context.Context, u model.User) (model.User, error) {
	var now = time.Now().UTC()
	var user model.User

	err := r.db.QueryRowContext(
		ctx,
		createUser,
		u.Username,
		u.PasswordHash,
		u.EmailEncrypted,
		now,
		now,
		u.Role,
	).Scan(
		&user.Id,
		&user.Username,
		&user.EmailEncrypted,
		&user.CreatedAt,
		&user.ModifiedAt,
		&user.Role,
	)
	if err != nil {
		return model.User{}, DBError(err)
	}

	return user, nil
}

func (r *Repository) GetUser(ctx context.Context, id int64) (model.User, error) {
	var user model.User

	err := r.db.QueryRowContext(
		ctx,
		getUser,
		id,
	).Scan(
		&user.Id,
		&user.Username,
		&user.EmailEncrypted,
		&user.CreatedAt,
		&user.ModifiedAt,
		&user.Role,
	)
	if err != nil {
		return model.User{}, DBError(err)
	}

	return user, nil
}

func (r *Repository) ListUsers(ctx context.Context) ([]model.User, error) {
	rows, err := r.db.QueryContext(
		ctx,
		listUsers,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var u model.User
		if err := rows.Scan(
			&u.Id,
			&u.Username,
			&u.EmailEncrypted,
			&u.CreatedAt,
			&u.ModifiedAt,
			&u.Role,
		); err != nil {
			return nil, DBError(err)
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *Repository) DeleteUser(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(
		ctx,
		deleteUser,
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
		return ErrUserNotFound
	}

	return nil
}
