package handlers_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Guldana11/gophermart/handlers"
	"github.com/Guldana11/gophermart/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockUserService struct {
	RegisterFunc func(ctx context.Context, login, password string) (*models.User, error)
	LoginFunc    func(ctx context.Context, login, password string) (*models.User, error)
}

func (m *mockUserService) Register(ctx context.Context, login, password string) (*models.User, error) {
	return m.RegisterFunc(ctx, login, password)
}

func (m *mockUserService) Login(ctx context.Context, login, password string) (*models.User, error) {
	return m.LoginFunc(ctx, login, password)
}

func TestRegisterHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		body           string
		contentType    string
		mockRegister   func(ctx context.Context, login, password string) (*models.User, error)
		expectedStatus int
	}{
		{
			name:        "success",
			body:        `{"login":"user","password":"pass"}`,
			contentType: "application/json",
			mockRegister: func(ctx context.Context, login, password string) (*models.User, error) {
				return &models.User{ID: "1", Login: login}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "extra field",
			body:           `{"login":"user","password":"pass","x":1}`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid json",
			body:           `{"login":"user",,,,}`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty login",
			body:           `{"login":"","password":"pass"}`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty password",
			body:           `{"login":"user","password":""}`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "both empty",
			body:           `{"login":"","password":""}`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "login not string",
			body:           `{"login":123,"password":"pass"}`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "password not string",
			body:           `{"login":"user","password":123}`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty body",
			body:           ``,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "null body",
			body:           `null`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid content-type",
			body:           `{"login":"user","password":"pass"}`,
			contentType:    "text/plain",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "too large payload (>1MB)",
			body:           `{"login":"` + strings.Repeat("A", 2_000_000) + `","password":"pass"}`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "login already exists",
			body:        `{"login":"user","password":"pass"}`,
			contentType: "application/json",
			mockRegister: func(ctx context.Context, login, password string) (*models.User, error) {
				return nil, errors.New("login already exists")
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:        "internal server error",
			body:        `{"login":"user","password":"pass"}`,
			contentType: "application/json",
			mockRegister: func(ctx context.Context, login, password string) (*models.User, error) {
				return nil, errors.New("db error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockUserService{
				RegisterFunc: func(ctx context.Context, login, password string) (*models.User, error) {
					if tt.mockRegister != nil {
						return tt.mockRegister(ctx, login, password)
					}
					return &models.User{ID: "1", Login: login}, nil
				},
			}

			r := gin.New()
			r.POST("/api/user/register", handlers.RegisterHandler(mockSvc))

			req := httptest.NewRequest("POST", "/api/user/register", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", tt.contentType)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestLoginHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		body           string
		contentType    string
		mockLogin      func(ctx context.Context, login, password string) (*models.User, error)
		expectedStatus int
	}{
		{
			name:        "success",
			body:        `{"login":"user","password":"pass"}`,
			contentType: "application/json",
			mockLogin: func(ctx context.Context, login, password string) (*models.User, error) {
				return &models.User{ID: "1", Login: login}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "invalid login or password",
			body:        `{"login":"user","password":"pass"}`,
			contentType: "application/json",
			mockLogin: func(ctx context.Context, login, password string) (*models.User, error) {
				return nil, errors.New("invalid login or password")
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "empty login",
			body:           `{"login":"","password":"pass"}`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty password",
			body:           `{"login":"user","password":""}`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "both empty",
			body:           `{"login":"","password":""}`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "login not string",
			body:           `{"login":123,"password":"pass"}`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "password not string",
			body:           `{"login":"user","password":123}`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "extra field",
			body:           `{"login":"user","password":"pass","x":1}`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid json",
			body:           `{"login":"user",,,,}`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty body",
			body:           ``,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "null body",
			body:           `null`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid content-type",
			body:           `{"login":"user","password":"pass"}`,
			contentType:    "text/plain",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "too large payload (>1MB)",
			body:           `{"login":"` + strings.Repeat("A", 2_000_000) + `","password":"pass"}`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "internal server error",
			body:        `{"login":"user","password":"pass"}`,
			contentType: "application/json",
			mockLogin: func(ctx context.Context, login, password string) (*models.User, error) {
				return nil, errors.New("db error")
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockUserService{
				LoginFunc: func(ctx context.Context, login, password string) (*models.User, error) {
					if tt.mockLogin != nil {
						return tt.mockLogin(ctx, login, password)
					}
					return &models.User{ID: "1", Login: login}, nil
				},
			}

			r := gin.New()
			r.POST("/api/user/login", handlers.LoginHandler(mockSvc))

			req := httptest.NewRequest("POST", "/api/user/login", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", tt.contentType)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
