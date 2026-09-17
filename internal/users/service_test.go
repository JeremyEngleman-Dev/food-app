package users

import (
	"context"
	"testing"

	m "foodapp/internal/models"
)

// Mock Repo
type mockRepo struct {
	mockCreateUser         func(ctx context.Context, u m.User) (m.User, error)
	mockGetUserById        func(ctx context.Context, id int64) (m.User, error)
	mockGetUserByEmailHash func(ctx context.Context, emailHash string) (m.User, error)
	mockListUsers          func(ctx context.Context) ([]m.User, error)
	mockDeleteUser         func(ctx context.Context, id int64) error
}

func (r *mockRepo) CreateUser(ctx context.Context, u m.User) (m.User, error) {
	return r.mockCreateUser(ctx, u)
}

func (r *mockRepo) GetUserById(ctx context.Context, id int64) (m.User, error) {
	return r.mockGetUserById(ctx, id)
}

func (r *mockRepo) GetUserByEmailHash(ctx context.Context, emailHash string) (m.User, error) {
	return r.mockGetUserByEmailHash(ctx, emailHash)
}

func (r *mockRepo) ListUsers(ctx context.Context) ([]m.User, error) {
	return r.mockListUsers(ctx)
}

func (r *mockRepo) DeleteUser(ctx context.Context, id int64) error {
	return r.mockDeleteUser(ctx, id)
}

// Mock Hash
type mockHash struct {
	mockHashPassword func(data string) (string, error)
	mockHashEmail    func(email string) string
}

func (h *mockHash) HashPassword(data string) (string, error) {
	return h.mockHashPassword(data)
}

func (h *mockHash) HashEmail(email string) string {
	return h.mockHashEmail(email)
}

func (h *mockHash) ValidateHashPassword(hash, data string) bool {
	return false
}

func (h *mockHash) HashRefreshToken(token []byte) string {
	return ""
}

// Mock Encryption
type mockEncryption struct {
	mockEncrypt func(value string) (string, error)
}

func (e *mockEncryption) Encrypt(value string) (string, error) {
	return e.mockEncrypt(value)
}

func (e *mockEncryption) Decrypt(value string) (string, error) {
	return "", nil
}

// Tests
func Test_Service_CreateUser(t *testing.T) {
	ctx := context.Background()

	t.Run("Successful_creation", func(t *testing.T) {
		var got m.User
		svc := NewService(
			&mockRepo{mockCreateUser: func(_ context.Context, u m.User) (m.User, error) {
				got = u
				return m.User{Id: 1}, nil
			}},
			&mockEncryption{mockEncrypt: func(value string) (string, error) { return "ENCRYPTED:" + value, nil }},
			&mockHash{
				mockHashPassword: func(data string) (string, error) { return "HASHEDPW:" + data, nil },
				mockHashEmail:    func(email string) string { return "HASHEDEMAIL:" + email },
			},
		)

		user, err := svc.CreateUser(ctx, m.CreateUser{
			DisplayName: "TestUser",
			Password:    "password",
			Email:       "TestUser@Test.com",
			Role:        "admin",
		})
		if err != nil {
			t.Fatalf("Test user creation error: %v", err)
		}
		if user.Id != 1 {
			t.Errorf("Test user expected 1, got [%d]", user.Id)
		}

		want := m.User{
			DisplayName:    "TestUser",
			PasswordHash:   "HASHEDPW:password",
			EmailHash:      "HASHEDEMAIL:TestUser@Test.com",
			EmailEncrypted: "ENCRYPTED:TestUser@Test.com",
			Role:           "admin",
		}
		if got != want {
			t.Errorf("Repo received %+v, want %+v", got, want)
		}
	})
}
