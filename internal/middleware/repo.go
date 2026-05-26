package middleware

import (
	m "Feastio/internal/models"
	"Feastio/internal/platform/database"
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

var (
	checkLoginUser = `SELECT "sessionId", "userId", "createdAt", "expiresAt", "role"
	FROM login_sessions WHERE "sessionId"=$1`
)

func (r *Repository) GetSessionBySessionId(ctx context.Context, sessionId string) (m.Session, error) {
	var session m.Session

	err := r.db.QueryRowContext(
		ctx,
		checkLoginUser,
		sessionId,
	).Scan(
		&session.SessionId,
		&session.UserId,
		&session.CreatedAt,
		&session.ExpiresAt,
		&session.Role,
	)
	if err != nil {
		return m.Session{}, database.DBError(err)
	}

	return session, nil
}
