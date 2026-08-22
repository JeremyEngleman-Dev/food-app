package auth

import (
	"context"
	m "foodapp/internal/models"
)

type UserProvider interface {
	GetUserById(ctx context.Context, id int64) (m.User, error)
}

type SessionProvider interface {
	GetSessionBySessionId(ctx context.Context, sessionId string) (m.Session, error)
}

type Service struct {
	users    UserProvider
	sessions SessionProvider
}

func NewService(users UserProvider, sessions SessionProvider) *Service {
	return &Service{
		users:    users,
		sessions: sessions,
	}
}

func (s *Service) ValidateSession(ctx context.Context, sessionId string) (m.UserContext, error) {
	session, err := s.sessions.GetSessionBySessionId(ctx, sessionId)
	if err != nil {
		return m.UserContext{}, err
	}

	user, err := s.users.GetUserById(ctx, session.UserId)
	if err != nil {
		return m.UserContext{}, err
	}

	return m.UserContext{
		UserId: session.UserId,
		Role:   user.Role,
	}, nil
}
