package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	m "foodapp/internal/models"
	"foodapp/internal/platform/database"
)

type Repository interface {
	CreateRefreshToken(ctx context.Context, token m.RefreshToken) error
	GetRefreshTokenByTokenHash(ctx context.Context, token string) (m.RefreshToken, error)
	RefreshTokenExists(ctx context.Context, userId int64) (bool, error)
	UpdateRefreshByUserId(ctx context.Context, userId int64, req m.UpdateRefreshToken) error
	DeleteRefreshTokenByUserId(ctx context.Context, userId int64) error
}

type repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) Repository {
	return &repo{db: db}
}

var (
	createRefreshToken = `INSERT INTO jwt (
		"refreshTokenHash", "userId", "createdAt", "expiresAt", "revokedAt"
	) VALUES ($1, $2, $3, $4, $5)`
	getRefreshTokenByTokenHash = `SELECT "refreshTokenHash", "userId", "createdAt", "expiresAt", "revokedAt"
		FROM jwt WHERE "refreshTokenHash" = $1`
	refreshTokenExists         = `SELECT 1 FROM jwt WHERE "userId" = $1`
	deleteRefreshTokenByUserId = `DELETE FROM jwt WHERE "userId" = $1`
)

func (r *repo) CreateRefreshToken(ctx context.Context, token m.RefreshToken) error {
	_, err := r.db.ExecContext(
		ctx,
		createRefreshToken,
		token.RefreshTokenHash,
		token.UserId,
		token.CreatedAt,
		token.ExpiresAt,
		token.RevokedAt,
	)
	fmt.Println(err)
	if err != nil {
		return database.DBError(err)
	}

	return nil
}

func (r *repo) GetRefreshTokenByTokenHash(ctx context.Context, token string) (m.RefreshToken, error) {
	var refreshToken m.RefreshToken

	err := r.db.QueryRowContext(
		ctx,
		getRefreshTokenByTokenHash,
		token,
	).Scan(
		&refreshToken.RefreshTokenHash,
		&refreshToken.UserId,
		&refreshToken.CreatedAt,
		&refreshToken.ExpiresAt,
		&refreshToken.RevokedAt,
	)
	if err != nil {
		return m.RefreshToken{}, database.DBError(err)
	}

	return refreshToken, nil
}

func (r *repo) RefreshTokenExists(ctx context.Context, userId int64) (bool, error) {
	var exists int

	err := r.db.QueryRowContext(
		ctx,
		refreshTokenExists,
		userId,
	).Scan(&exists)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	return true, nil
}

func (r *repo) UpdateRefreshByUserId(ctx context.Context, userId int64, req m.UpdateRefreshToken) error {
	updates := database.PatchMapUpdates(req)

	query, args := database.PatchQueryBuilder("jwt", updates, "userId", userId)

	_, err := r.db.ExecContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		fmt.Println(err)
		return database.DBError(err)
	}

	return nil
}

func (r *repo) DeleteRefreshTokenByUserId(ctx context.Context, userId int64) error {
	_, err := r.db.ExecContext(
		ctx,
		deleteRefreshTokenByUserId,
		userId,
	)
	if err != nil {
		return database.DBError(err)
	}

	return nil
}
