package service

import (
	"context"
	"errors"
	"log"

	"Feastio/internal/dto"
	"Feastio/internal/model"
	"Feastio/internal/repository"
)

// Error Definitions
var (
	ErrUserNotFound      = errors.New("User not found")
	ErrDuplicateUsername = errors.New("Username already exists")
	ErrDuplicateEmail    = errors.New("Email already exists")
)

// Functions
func (s *Service) CreateUser(ctx context.Context, u dto.CreateUser) (dto.DisplayUser, error) {
	hashedPassword, err := HashPassword(u.Password)
	if err != nil {
		return dto.DisplayUser{}, err
	}
	encryptedEmail, err := EncryptData(u.Email, s.emailKey)
	if err != nil {
		return dto.DisplayUser{}, err
	}

	userCreation := model.User{
		Username:       u.Username,
		PasswordHash:   hashedPassword,
		EmailEncrypted: encryptedEmail,
		Role:           u.Role,
	}

	user, err := s.repo.CreateUser(ctx, userCreation)

	if err != nil {
		if errors.Is(err, repository.ErrDuplicateUsername) {
			return dto.DisplayUser{}, ErrDuplicateUsername
		}
		if errors.Is(err, repository.ErrDuplicateEmail) {
			return dto.DisplayUser{}, ErrDuplicateEmail
		}
		if errors.Is(err, repository.ErrNotFound) {
			return dto.DisplayUser{}, ErrUserNotFound
		}
		return dto.DisplayUser{}, err
	}

	response, err := s.DisplayUser(user)
	if err != nil {
		return dto.DisplayUser{}, err
	}

	return response, nil
}

func (s *Service) GetUser(ctx context.Context, id int64) (dto.DisplayUser, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return dto.DisplayUser{}, ErrUserNotFound
		}
		return dto.DisplayUser{}, err
	}

	response, err := s.DisplayUser(user)
	if err != nil {
		return dto.DisplayUser{}, err
	}

	return response, nil
}

func (s *Service) ListUsers(ctx context.Context) ([]dto.DisplayUser, error) {
	users, err := s.repo.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]dto.DisplayUser, len(users))
	for i, u := range users {
		response[i], err = s.DisplayUser(u)
		if err != nil {
			// Handle this better, shame
			log.Println("Model to dto conversion error")
		}
	}

	return response, nil
}

func (s *Service) DeleteUser(ctx context.Context, id int64) error {
	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrIngredientNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	return nil
}

// Support Functions
func (s *Service) DisplayUser(m model.User) (dto.DisplayUser, error) {
	email, err := DecryptData(m.EmailEncrypted, s.emailKey)
	if err != nil {
		return dto.DisplayUser{}, err
	}

	user := dto.DisplayUser{
		Id:         m.Id,
		Username:   m.Username,
		Email:      email,
		CreatedAt:  m.CreatedAt,
		ModifiedAt: m.ModifiedAt,
		Role:       m.Role,
	}

	return user, nil
}
