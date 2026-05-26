package middleware

import (
	m "Feastio/internal/models"
	"context"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) ValidateSession(ctx context.Context, sessionId string) (m.UserContext, error) {
	session, err := s.repo.GetSessionBySessionId(ctx, sessionId)
	if err != nil {
		return m.UserContext{}, err
	}

	return m.UserContext{
		UserId: session.UserId,
		Role:   session.Role,
	}, nil
}
