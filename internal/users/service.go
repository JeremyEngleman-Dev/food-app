package users

import (
	"context"
	"errors"
	m "foodapp/internal/models"
	"foodapp/internal/platform/database"
	"foodapp/internal/platform/security"
)

type Service struct {
	repo       Repository
	encryption security.Encryption
	hash       security.Hash
}

func NewService(repo Repository, encryption security.Encryption, hash security.Hash) *Service {
	return &Service{
		repo:       repo,
		encryption: encryption,
		hash:       hash,
	}
}

// Functions
func (s *Service) CreateUser(ctx context.Context, u m.CreateUser) (m.User, error) {
	hashedEmail := s.hash.HashEmail(u.Email)

	hashedPassword, err := s.hash.HashPassword(u.Password)
	if err != nil {
		return m.User{}, err
	}

	encryptedEmail, err := s.encryption.Encrypt(u.Email)
	if err != nil {
		return m.User{}, err
	}

	userCreation := m.User{
		DisplayName:    u.DisplayName,
		PasswordHash:   hashedPassword,
		EmailHash:      hashedEmail,
		EmailEncrypted: encryptedEmail,
		Role:           u.Role,
	}

	user, err := s.repo.CreateUser(ctx, userCreation)
	if err != nil {
		var appErr *database.AppError
		if errors.As(err, &appErr) {
			if appErr.Type == database.ErrTypeNotFound {
				return m.User{}, &database.AppError{Type: database.ErrTypeFailedCreation, Err: err}
			}
		}
		return m.User{}, err
	}

	return user, nil
}

func (s *Service) GetUserById(ctx context.Context, id int64) (m.User, error) {
	user, err := s.repo.GetUserById(ctx, id)
	if err != nil {
		return m.User{}, err
	}

	return user, nil
}

func (s *Service) GetUserByEmailHash(ctx context.Context, emailHash string) (m.User, error) {
	user, err := s.repo.GetUserByEmailHash(ctx, emailHash)
	if err != nil {
		return m.User{}, err
	}

	return user, nil
}

func (s *Service) ListUsers(ctx context.Context) ([]m.User, error) {
	users, err := s.repo.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Service) DeleteUserById(ctx context.Context, id int64) error {
	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
