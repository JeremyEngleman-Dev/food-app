package users

import (
	"Feastio/internal/platform/database"
	"context"
	"database/sql"
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
	createUser = `
	INSERT INTO users (
		"displayName", "passwordHash", "emailHash", "emailEncrypted",
		"createdAt", "modifiedAt", "role"
	) VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING "id", "displayName", "emailEncrypted", "createdAt", "modifiedAt", "role"
	`
	getUser = `
	SELECT "id", "displayName", "emailEncrypted", "createdAt", "modifiedAt","role"
	FROM users WHERE id=$1
	`
	listUsers = `SELECT "id", "displayName", "emailEncrypted", "createdAt", "modifiedAt","role"
	FROM users`
	deleteUser      = `DELETE FROM users WHERE id = $1`
	checkEmailExist = `SELECT "emailHash" FROM users WHERE "emailHash" = $1`
)

// Functions
func (r *Repository) CreateUser(ctx context.Context, u User) (User, error) {
	var now = time.Now().UTC()
	var user User

	err := r.db.QueryRowContext(
		ctx,
		createUser,
		u.DisplayName,
		u.PasswordHash,
		u.EmailHash,
		u.EmailEncrypted,
		now,
		now,
		u.Role,
	).Scan(
		&user.Id,
		&user.DisplayName,
		&user.EmailEncrypted,
		&user.CreatedAt,
		&user.ModifiedAt,
		&user.Role,
	)
	if err != nil {
		return User{}, database.DBError(err)
	}

	return user, nil
}

func (r *Repository) GetUser(ctx context.Context, id int64) (User, error) {
	var user User

	err := r.db.QueryRowContext(
		ctx,
		getUser,
		id,
	).Scan(
		&user.Id,
		&user.DisplayName,
		&user.EmailEncrypted,
		&user.CreatedAt,
		&user.ModifiedAt,
		&user.Role,
	)
	if err != nil {
		return User{}, database.DBError(err)
	}

	return user, nil
}

func (r *Repository) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := r.db.QueryContext(
		ctx,
		listUsers,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
		if err := rows.Scan(
			&u.Id,
			&u.DisplayName,
			&u.EmailEncrypted,
			&u.CreatedAt,
			&u.ModifiedAt,
			&u.Role,
		); err != nil {
			return nil, database.DBError(err)
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
		return database.DBError(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return database.ErrNotFound
	}

	return nil
}
