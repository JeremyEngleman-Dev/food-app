package mappings

import (
	m "foodapp/internal/models"
	"foodapp/internal/platform/security"
)

type UserMapping interface {
	ToDisplayUser(user m.User) (m.DisplayUser, error)
}

type userMapping struct {
	encryption security.Encryption
}

func NewUserMap(encryption security.Encryption) *userMapping {
	return &userMapping{
		encryption: encryption,
	}
}

func (u *userMapping) ToDisplayUser(user m.User) (m.DisplayUser, error) {
	email, err := u.encryption.Decrypt(user.EmailEncrypted)
	if err != nil {
		return m.DisplayUser{}, err
	}
	return m.DisplayUser{
		Id:          user.Id,
		DisplayName: user.DisplayName,
		Email:       email,
		CreatedAt:   user.CreatedAt,
		ModifiedAt:  user.ModifiedAt,
		Role:        user.Role,
	}, nil
}
