package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Guldana11/gophermart/models"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct {
	CreateUserFunc     func(ctx context.Context, login, password string) (*models.User, error)
	GetUserByLoginFunc func(ctx context.Context, login string) (*models.User, error)
}

func (m *mockUserRepo) CreateUser(ctx context.Context, login, password string) (*models.User, error) {
	return m.CreateUserFunc(ctx, login, password)
}

func (m *mockUserRepo) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	return m.GetUserByLoginFunc(ctx, login)
}

func (m *mockUserRepo) GetUserPoints(ctx context.Context, userID string) (float64, float64, error) {
	return 100, 0, nil
}

func (m *mockUserRepo) WithdrawPoints(ctx context.Context, userID string, order string, sum float64) (float64, error) {
	return 100 - sum, nil
}
func TestCheckPasswordHash(t *testing.T) {
	password := "secret"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	tests := []struct {
		name string
		pass string
		hash string
		want bool
	}{
		{"correct password", password, string(hash), true},
		{"wrong password", "wrong", string(hash), false},
		{"empty password", "", string(hash), false},
		{"empty hash", password, "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, CheckPasswordHash(tt.pass, tt.hash))
		})
	}
}

func TestNewUserService(t *testing.T) {
	mockRepo := &mockUserRepo{}
	svc := NewUserService(mockRepo)
	assert.NotNil(t, svc)
	assert.Equal(t, mockRepo, svc.repo)
}

func TestUserService_Register(t *testing.T) {
	tests := []struct {
		name    string
		login   string
		pass    string
		mockFn  func(ctx context.Context, login, pass string) (*models.User, error)
		wantErr bool
	}{
		{"success", "user", "pass", func(ctx context.Context, l, p string) (*models.User, error) {
			return &models.User{ID: "1", Login: l}, nil
		}, false},
		{"empty login", "", "pass", nil, true},
		{"empty password", "user", "", nil, true},
		{"db error", "user", "pass", func(ctx context.Context, l, p string) (*models.User, error) {
			return nil, errors.New("db error")
		}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockUserRepo{
				CreateUserFunc: func(ctx context.Context, login, password string) (*models.User, error) {
					if tt.mockFn == nil {
						return &models.User{ID: "1", Login: login}, nil
					}
					return tt.mockFn(ctx, login, password)
				},
			}
			svc := NewUserService(mockRepo)
			got, err := svc.Register(context.Background(), tt.login, tt.pass)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.login, got.Login)
			}
		})
	}
}

func TestUserService_Authenticate(t *testing.T) {
	pass := "secret"
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	tests := []struct {
		name    string
		login   string
		pass    string
		mockFn  func(ctx context.Context, login string) (*models.User, error)
		wantErr bool
	}{
		{"success", "user", pass, func(ctx context.Context, l string) (*models.User, error) {
			return &models.User{ID: "1", Login: l, PasswordHash: string(hash)}, nil
		}, false},
		{"wrong password", "user", "wrong", func(ctx context.Context, l string) (*models.User, error) {
			return &models.User{ID: "1", Login: l, PasswordHash: string(hash)}, nil
		}, true},
		{"user not found", "user", pass, func(ctx context.Context, l string) (*models.User, error) {
			return nil, errors.New("not found")
		}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockUserRepo{
				GetUserByLoginFunc: tt.mockFn,
			}
			svc := NewUserService(mockRepo)
			got, err := svc.Authenticate(context.Background(), tt.login, tt.pass)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.login, got.Login)
			}
		})
	}
}

func TestUserService_Login(t *testing.T) {
	pass := "secret"
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	tests := []struct {
		name    string
		login   string
		pass    string
		mockFn  func(ctx context.Context, login string) (*models.User, error)
		wantErr bool
	}{
		{"success", "user", pass, func(ctx context.Context, l string) (*models.User, error) {
			return &models.User{ID: "1", Login: l, PasswordHash: string(hash)}, nil
		}, false},
		{"wrong password", "user", "wrong", func(ctx context.Context, l string) (*models.User, error) {
			return &models.User{ID: "1", Login: l, PasswordHash: string(hash)}, nil
		}, true},
		{"user not found", "user", pass, func(ctx context.Context, l string) (*models.User, error) {
			return nil, errors.New("not found")
		}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockUserRepo{
				GetUserByLoginFunc: tt.mockFn,
			}
			svc := NewUserService(mockRepo)
			got, err := svc.Login(context.Background(), tt.login, tt.pass)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.login, got.Login)
			}
		})
	}
}
