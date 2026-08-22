package sessions

import (
	"context"
	"database/sql"
	"fmt"
	m "foodapp/internal/models"
	"foodapp/internal/platform/database"
)

type Repository interface {
	CreateSession(ctx context.Context, s m.Session) error
	GetSessionBySessionId(ctx context.Context, sessionId string) (m.Session, error)
	DeleteSessionById(ctx context.Context, sessionId string) error
}

type repo struct {
	db *sql.DB
}

func Newrepo(db *sql.DB) Repository {
	return &repo{db: db}
}

// Database queries
var (
	loginUser = `INSERT INTO login_sessions (
		"sessionId", "userId", "createdAt", "expiresAt"
	) VALUES ($1, $2, $3, $4)`
	checkLoginUser = `SELECT "sessionId", "userId", "createdAt", "expiresAt"
	FROM login_sessions WHERE "sessionId"=$1`
	logoutUser = `DELETE FROM login_sessions WHERE "sessionId"=$1`
)

func (r *repo) CreateSession(ctx context.Context, s m.Session) error {
	_, err := r.db.ExecContext(
		ctx,
		loginUser,
		s.SessionId,
		s.UserId,
		s.CreatedAt,
		s.ExpiresAt,
	)

	if err != nil {
		fmt.Println(err)
		return database.DBError(err)
	}

	return nil
}

func (r *repo) GetSessionBySessionId(ctx context.Context, sessionId string) (m.Session, error) {
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
	)
	if err != nil {
		return m.Session{}, database.DBError(err)
	}

	return session, nil
}

func (r *repo) DeleteSessionById(ctx context.Context, sessionId string) error {
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
