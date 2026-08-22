package sessions

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	m "foodapp/internal/models"
	"foodapp/internal/platform/security"
	"time"
)

type UserProvider interface {
	GetUserByEmailHash(ctx context.Context, emailHash string) (m.User, error)
}

type Service struct {
	repo  Repository
	cfg   security.SecurityKeys
	users UserProvider
}

func NewService(repo Repository, cfg security.SecurityKeys, usersSvc UserProvider) *Service {
	return &Service{
		repo:  repo,
		cfg:   cfg,
		users: usersSvc,
	}
}

var (
	ErrBadCredentials = errors.New("Email or password incorrect")
)

func (s *Service) LoginUser(ctx context.Context, login m.LoginInfo) (m.Session, error) {
	loginEmail := security.HashData(login.Email, s.cfg.Hash)

	user, err := s.users.GetUserByEmailHash(ctx, loginEmail)
	if err != nil {
		return m.Session{}, err
	}

	checkUserPassword := security.CheckHashPassword(user.PasswordHash, login.Password)
	if !checkUserPassword {
		return m.Session{}, ErrBadCredentials
	}

	sessionId, err := generateSessionID()
	if err != nil {
		return m.Session{}, err
	}
	now := time.Now().UTC()

	session := m.Session{
		SessionId: sessionId,
		UserId:    user.Id,
		CreatedAt: now,
		ExpiresAt: now.Add(15 * time.Minute),
	}

	err = s.repo.CreateSession(ctx, session)
	if err != nil {
		return m.Session{}, err
	}

	return session, nil
}

func (s *Service) LogoutUser(ctx context.Context, sessionId string) error {
	err := s.repo.DeleteSessionById(ctx, sessionId)

	return err
}

// Support Functions
func generateSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
