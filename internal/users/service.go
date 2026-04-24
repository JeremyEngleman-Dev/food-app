package users

import (
	"context"
	"errors"
	"log"

	"Feastio/internal/platform/database"
)

type Service struct {
	repo *Repository
	cfg  UserKeys
}

func NewService(repo *Repository, cfg UserKeys) *Service {
	return &Service{
		repo: repo,
		cfg:  cfg,
	}
}

// Error Definitions
var (
	ErrUserNotFound   = errors.New("User not found")
	ErrDuplicateEmail = errors.New("Email already exists")
)

// Functions
func (s *Service) CreateUser(ctx context.Context, u CreateUser) (DisplayUser, error) {
	hashedEmail := HashData(u.Email, s.cfg.hash)

	hashedPassword, err := HashPassword(u.Password)
	if err != nil {
		return DisplayUser{}, err
	}
	encryptedEmail, err := EncryptData(u.Email, s.cfg.encryption)
	if err != nil {
		return DisplayUser{}, err
	}

	userCreation := User{
		DisplayName:    u.DisplayName,
		PasswordHash:   hashedPassword,
		EmailHash:      hashedEmail,
		EmailEncrypted: encryptedEmail,
		Role:           u.Role,
	}

	user, err := s.repo.CreateUser(ctx, userCreation)

	if err != nil {
		if errors.Is(err, database.ErrDuplicateEmail) {
			return DisplayUser{}, ErrDuplicateEmail
		}
		if errors.Is(err, database.ErrNotFound) {
			return DisplayUser{}, ErrUserNotFound
		}
		return DisplayUser{}, err
	}

	response, err := s.DisplayUser(user)
	if err != nil {
		return DisplayUser{}, err
	}

	return response, nil
}

func (s *Service) GetUser(ctx context.Context, id int64) (DisplayUser, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return DisplayUser{}, ErrUserNotFound
		}
		return DisplayUser{}, err
	}

	response, err := s.DisplayUser(user)
	if err != nil {
		return DisplayUser{}, err
	}

	return response, nil
}

func (s *Service) ListUsers(ctx context.Context) ([]DisplayUser, error) {
	users, err := s.repo.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]DisplayUser, len(users))
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
		if errors.Is(err, database.ErrNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	return nil
}

// Support Functions
func (s *Service) DisplayUser(m User) (DisplayUser, error) {
	email, err := DecryptData(m.EmailEncrypted, s.cfg.encryption)
	if err != nil {
		return DisplayUser{}, err
	}

	user := DisplayUser{
		Id:          m.Id,
		DisplayName: m.DisplayName,
		Email:       email,
		CreatedAt:   m.CreatedAt,
		ModifiedAt:  m.ModifiedAt,
		Role:        m.Role,
	}

	return user, nil
}
