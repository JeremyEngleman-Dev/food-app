package sessions

import (
	"context"
	"database/sql"
	"fmt"
	m "foodapp/internal/models"
	"foodapp/internal/platform/database"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Database queries
var (
	loginUser = `INSERT INTO login_sessions (
		"sessionId", "userId", "createdAt", "expiresAt", "role"
	) VALUES ($1, $2, $3, $4, $5)`
	logoutUser = `DELETE FROM login_sessions WHERE "sessionId"=$1`
)

func (r *Repository) CreateLoginSession(ctx context.Context, s m.Session) error {
	_, err := r.db.ExecContext(
		ctx,
		loginUser,
		s.SessionId,
		s.UserId,
		s.CreatedAt,
		s.ExpiresAt,
		s.Role,
	)

	if err != nil {
		fmt.Println(err)
		return database.DBError(err)
	}

	return nil
}

func (r *Repository) DeleteLoginSession(ctx context.Context, sessionId string) error {
	_, err := r.db.ExecContext(
		ctx,
		logoutUser,
		sessionId,
	)
	if err != nil {
		return database.DBError(err)
	}

	return nil
}
