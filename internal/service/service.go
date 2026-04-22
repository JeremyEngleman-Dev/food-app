package service

import (
	"Feastio/internal/repository"
	"crypto/cipher"
)

type Service struct {
	repo     *repository.Repository
	emailKey cipher.Block
}

func NewService(repo *repository.Repository, emailKey cipher.Block) *Service {
	return &Service{
		repo:     repo,
		emailKey: emailKey,
	}
}
