package users

import (
	"context"
	"errors"
	m "foodapp/internal/models"
	"foodapp/internal/platform/database"
	"foodapp/internal/platform/security"
	"log"
)

type Service struct {
	repo *Repository
	cfg  security.SecurityKeys
}

func NewService(repo *Repository, cfg security.SecurityKeys) *Service {
	return &Service{
		repo: repo,
		cfg:  cfg,
	}
}

// Functions
func (s *Service) CreateUser(ctx context.Context, u m.CreateUser) (m.DisplayUser, error) {
	hashedEmail := security.HashData(u.Email, s.cfg.Hash)

	hashedPassword, err := security.HashPassword(u.Password)
	if err != nil {
		return m.DisplayUser{}, err
	}
	encryptedEmail, err := security.EncryptData(u.Email, s.cfg.Encryption)
	if err != nil {
		return m.DisplayUser{}, err
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
				return m.DisplayUser{}, &database.AppError{Type: database.ErrTypeFailedCreation, Err: err}
			}
		}
		return m.DisplayUser{}, err
	}

	response, err := s.DisplayUser(user)
	if err != nil {
		return m.DisplayUser{}, err
	}

	return response, nil
}

func (s *Service) GetUserById(ctx context.Context, id int64) (m.DisplayUser, error) {
	user, err := s.repo.GetUserById(ctx, id)
	if err != nil {
		return m.DisplayUser{}, err
	}

	response, err := s.DisplayUser(user)
	if err != nil {
		return m.DisplayUser{}, err
	}

	return response, nil
}

func (s *Service) GetUserByEmailHash(ctx context.Context, emailHash string) (m.User, error) {
	user, err := s.repo.GetUserByEmailHash(ctx, emailHash)
	if err != nil {
		return m.User{}, err
	}

	return user, nil
}

func (s *Service) ListUsers(ctx context.Context) ([]m.DisplayUser, error) {
	users, err := s.repo.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]m.DisplayUser, len(users))
	for i, u := range users {
		response[i], err = s.DisplayUser(u)
		if err != nil {
			// Handle this better, shame
			log.Println("Model to dto conversion error")
		}
	}

	return response, nil
}

func (s *Service) DeleteUserById(ctx context.Context, id int64) error {
	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// Support Functions
func (s *Service) DisplayUser(u m.User) (m.DisplayUser, error) {
	email, err := security.DecryptData(u.EmailEncrypted, s.cfg.Encryption)
	if err != nil {
		return m.DisplayUser{}, err
	}

	user := m.DisplayUser{
		Id:          u.Id,
		DisplayName: u.DisplayName,
		Email:       email,
		CreatedAt:   u.CreatedAt,
		ModifiedAt:  u.ModifiedAt,
		Role:        u.Role,
	}

	return user, nil
}
