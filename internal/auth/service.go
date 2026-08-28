package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	m "foodapp/internal/models"
	"foodapp/internal/platform/security"
)

type UserProvider interface {
	GetUserByEmailHash(ctx context.Context, emailHash string) (m.User, error)
	GetUserById(ctx context.Context, id int64) (m.User, error)
}

type Service struct {
	repo  Repository
	hash  security.Hash
	token security.Token
	users UserProvider
}

func NewService(repo Repository, hash security.Hash, usersSvc UserProvider, token security.Token) *Service {
	return &Service{
		repo:  repo,
		hash:  hash,
		token: token,
		users: usersSvc,
	}
}

var (
	ErrBadCredentials      = errors.New("Email or password incorrect")
	ErrExpiredRefreshToken = errors.New("The refresh token has expired, new login required")
)

func (s *Service) LoginUser(ctx context.Context, login m.LoginInfo) (m.ResponseTokens, error) {
	loginEmail := s.hash.HashEmail(login.Email)

	user, err := s.users.GetUserByEmailHash(ctx, loginEmail)
	if err != nil {
		return m.ResponseTokens{}, ErrBadCredentials
	}

	checkUserPassword := s.hash.ValidateHashPassword(user.PasswordHash, login.Password)
	if !checkUserPassword {
		return m.ResponseTokens{}, ErrBadCredentials
	}

	refreshTokenExist, err := s.repo.RefreshTokenExists(ctx, user.Id)
	if err != nil {
		return m.ResponseTokens{}, err
	}

	token, err := s.token.CreateJwtToken(user)
	if err != nil {
		return m.ResponseTokens{}, err
	}

	refresh, err := s.token.CreateRefreshToken(30)
	if err != nil {
		return m.ResponseTokens{}, err
	}

	refreshEncoded := base64.StdEncoding.EncodeToString(refresh)
	refreshHashed := s.hash.HashRefreshToken(refresh)

	now := time.Now().UTC()
	expiresAt := now.AddDate(0, 1, 0)

	if refreshTokenExist {
		updatedToken := m.UpdateRefreshToken{
			RefreshTokenHash: &refreshHashed,
			CreatedAt:        &now,
			ExpiresAt:        &expiresAt,
		}

		err = s.repo.UpdateRefreshByUserId(ctx, user.Id, updatedToken)
		if err != nil {
			return m.ResponseTokens{}, err
		}
	} else {
		refreshToken := m.RefreshToken{
			RefreshTokenHash: refreshHashed,
			UserId:           user.Id,
			CreatedAt:        now,
			ExpiresAt:        expiresAt,
		}

		err = s.repo.CreateRefreshToken(
			ctx,
			refreshToken,
		)
		if err != nil {
			return m.ResponseTokens{}, err
		}
	}

	return m.ResponseTokens{
		Token:   token,
		Refresh: refreshEncoded,
	}, nil
}

func (s *Service) LogoutUser(ctx context.Context, userId int64) error {
	err := s.repo.DeleteRefreshTokenByUserId(ctx, userId)

	return err
}

func (s *Service) RefreshToken(ctx context.Context, refreshTokenEncoded string) (m.ResponseTokens, error) {
	refreshToken, err := base64.StdEncoding.DecodeString(refreshTokenEncoded)
	if err != nil {
		return m.ResponseTokens{}, err
	}

	tokenHashed := s.hash.HashRefreshToken([]byte(refreshToken))

	storedRefreshToken, err := s.repo.GetRefreshTokenByTokenHash(ctx, tokenHashed)
	if err != nil {
		return m.ResponseTokens{}, err
	}

	if storedRefreshToken.ExpiresAt.Before(time.Now()) {
		return m.ResponseTokens{}, ErrExpiredRefreshToken
	}

	user, err := s.users.GetUserById(ctx, storedRefreshToken.UserId)
	if err != nil {
		return m.ResponseTokens{}, err
	}

	newToken, err := s.token.CreateJwtToken(user)
	if err != nil {
		return m.ResponseTokens{}, err
	}

	refresh, err := s.token.CreateRefreshToken(30)
	if err != nil {
		return m.ResponseTokens{}, err
	}

	refreshEncoded := base64.StdEncoding.EncodeToString(refresh)
	refreshHashed := s.hash.HashRefreshToken(refresh)

	newExpiration := time.Now().UTC().AddDate(0, 1, 0)

	updateRefreshToken := m.UpdateRefreshToken{
		RefreshTokenHash: &refreshHashed,
		ExpiresAt:        &newExpiration,
	}

	err = s.repo.UpdateRefreshByUserId(ctx, user.Id, updateRefreshToken)
	if err != nil {
		fmt.Println("UpdateUser")
		return m.ResponseTokens{}, err
	}

	return m.ResponseTokens{
		Token:   newToken,
		Refresh: refreshEncoded,
	}, nil
}
