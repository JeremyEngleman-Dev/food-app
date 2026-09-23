package users

import (
	"context"
	"errors"
	"testing"

	m "foodapp/internal/models"
	"foodapp/internal/platform/database"
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

	t.Run("successful user creation", func(t *testing.T) {
		var repoReceived m.User
		svc := NewService(
			&mockRepo{mockCreateUser: func(_ context.Context, u m.User) (m.User, error) {
				repoReceived = u
				return m.User{Id: 1}, nil
			}},
			&mockEncryption{mockEncrypt: func(value string) (string, error) { return "ENCRYPTED:" + value, nil }},
			&mockHash{
				mockHashPassword: func(data string) (string, error) { return "HASHEDPW:" + data, nil },
				mockHashEmail:    func(email string) string { return "HASHEDEMAIL:" + email },
			},
		)

		got, err := svc.CreateUser(ctx, m.CreateUser{
			DisplayName: "TestUser",
			Password:    "password",
			Email:       "TestUser@Test.com",
			Role:        "admin",
		})
		if err != nil {
			t.Fatalf("svc.CreateUser() returned an unexpected error: %v", err)
		}

		repoWanted := m.User{
			DisplayName:    "TestUser",
			PasswordHash:   "HASHEDPW:password",
			EmailHash:      "HASHEDEMAIL:TestUser@Test.com",
			EmailEncrypted: "ENCRYPTED:TestUser@Test.com",
			Role:           "admin",
		}
		if repoReceived != repoWanted {
			t.Fatalf("repo received %+v, want %+v", repoReceived, repoWanted)
		}

		if got.Id != 1 {
			t.Errorf("svc.CreateUser() returned user id [%d], want 1", got.Id)
		}
	})

	t.Run("failed encryption", func(t *testing.T) {
		expectedError := errors.New("failed to encrypt")
		repoCalled := false

		svc := NewService(
			&mockRepo{mockCreateUser: func(_ context.Context, u m.User) (m.User, error) {
				repoCalled = true
				return m.User{}, nil
			}},
			&mockEncryption{mockEncrypt: func(value string) (string, error) { return "", expectedError }},
			&mockHash{
				mockHashPassword: func(data string) (string, error) { return data, nil },
				mockHashEmail:    func(email string) string { return email },
			},
		)

		user, err := svc.CreateUser(ctx, m.CreateUser{Password: "secretPassword", Email: "TestUser@Test.com"})
		if err != expectedError {
			t.Fatalf("svc.CreateUser() expected error %v, got %v", expectedError, err)
		}
		if user != (m.User{}) {
			t.Fatalf("Expected null user, got %v", user)
		}
		if repoCalled {
			t.Error("Repo should not have been called when encryption failed")
		}
	})

	t.Run("failed password hash", func(t *testing.T) {
		expectedError := errors.New("failed to hash")
		repoCalled := false

		svc := NewService(
			&mockRepo{mockCreateUser: func(_ context.Context, u m.User) (m.User, error) {
				repoCalled = true
				return m.User{}, nil
			}},
			&mockEncryption{mockEncrypt: func(value string) (string, error) { return value, nil }},
			&mockHash{
				mockHashPassword: func(data string) (string, error) { return "", expectedError },
				mockHashEmail:    func(email string) string { return email },
			},
		)

		user, err := svc.CreateUser(ctx, m.CreateUser{Password: "secretPassword", Email: "TestUser@Test.com"})
		if err != expectedError {
			t.Fatalf("svc.CreateUser() expected error %v, got %v", expectedError, err)
		}
		if user != (m.User{}) {
			t.Fatalf("Expected null user, got %v", user)
		}
		if repoCalled {
			t.Error("Repo should not have been called when password hashing failed")
		}
	})

	t.Run("unknown database error", func(t *testing.T) {
		expectedError := errors.New("Unknown database error")

		svc := NewService(
			&mockRepo{mockCreateUser: func(_ context.Context, u m.User) (m.User, error) {
				return m.User{}, &database.AppError{Type: database.ErrTypeDatabase, Err: expectedError}
			}},
			&mockEncryption{mockEncrypt: func(value string) (string, error) { return value, nil }},
			&mockHash{
				mockHashPassword: func(data string) (string, error) { return data, nil },
				mockHashEmail:    func(email string) string { return email },
			},
		)

		_, err := svc.CreateUser(ctx, m.CreateUser{Password: "secretPassword", Email: "TestUser@Test.com"})
		if err == nil {
			t.Fatalf("expected error, but no error was received")
		}
		var appErr *database.AppError
		if !errors.As(err, &appErr) {
			t.Fatalf("expected database.AppError, got %T", err)
		}
		if appErr.Type != database.ErrTypeDatabase {
			t.Errorf("expected type %v, got %v", database.ErrTypeDatabase, appErr.Type)
		}
	})
}

func Test_Service_GetUserById(t *testing.T) {
	ctx := context.Background()

	t.Run("successful user retrieval", func(t *testing.T) {
		svc := NewService(
			&mockRepo{mockGetUserById: func(_ context.Context, id int64) (m.User, error) {
				return m.User{
					Id:             1,
					DisplayName:    "TestUser",
					PasswordHash:   "HASHEDPW:secretPassword",
					EmailHash:      "HASHEDEMAIL:TestUser@Test.com",
					EmailEncrypted: "ENCRYPTED:TestUser@Test.com",
					Role:           "admin",
				}, nil
			}},
			&mockEncryption{mockEncrypt: func(value string) (string, error) { return value, nil }},
			&mockHash{
				mockHashPassword: func(data string) (string, error) { return data, nil },
				mockHashEmail:    func(email string) string { return email },
			},
		)

		user, err := svc.GetUserById(ctx, 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user.Id != 1 {
			t.Fatalf("got id [%d], want id [1]", user.Id)
		}

		expectedUser := m.User{
			Id:             1,
			DisplayName:    "TestUser",
			PasswordHash:   "HASHEDPW:secretPassword",
			EmailHash:      "HASHEDEMAIL:TestUser@Test.com",
			EmailEncrypted: "ENCRYPTED:TestUser@Test.com",
			Role:           "admin",
		}

		if user != expectedUser {
			t.Errorf("got user %v, want user %v", user, expectedUser)
		}
	})
}
