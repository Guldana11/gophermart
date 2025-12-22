package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Guldana11/gophermart/models"
	"github.com/Guldana11/gophermart/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockBalanceService struct {
	GetUserBalanceFunc func(ctx context.Context, userID string) (float64, float64, error)
	WithdrawFunc       func(ctx context.Context, userID, order string, sum float64) (float64, error)
	GetWithdrawalsFunc func(ctx context.Context, userID string) ([]models.WithdrawalResponse, error)
}

func (m *MockBalanceService) GetUserBalance(ctx context.Context, userID string) (float64, float64, error) {
	return m.GetUserBalanceFunc(ctx, userID)
}

func (m *MockBalanceService) Withdraw(ctx context.Context, userID, order string, sum float64) (float64, error) {
	return m.WithdrawFunc(ctx, userID, order, sum)
}

func (m *MockBalanceService) GetWithdrawals(ctx context.Context, userID string) ([]models.WithdrawalResponse, error) {
	return m.GetWithdrawalsFunc(ctx, userID)
}

func TestGetBalanceHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		userID       string
		mockFunc     func(ctx context.Context, userID string) (float64, float64, error)
		expectedCode int
		expectedBody string
	}{
		{
			name:   "unauthorized",
			userID: "",
			mockFunc: func(ctx context.Context, userID string) (float64, float64, error) {
				return 0, 0, nil
			},
			expectedCode: http.StatusUnauthorized,
			expectedBody: "",
		},
		{
			name:   "internal error",
			userID: "123",
			mockFunc: func(ctx context.Context, userID string) (float64, float64, error) {
				return 0, 0, errors.New("db error")
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error":"internal server error"}`,
		},
		{
			name:   "success",
			userID: "123",
			mockFunc: func(ctx context.Context, userID string) (float64, float64, error) {
				return 500.5, 42, nil
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"current":500.5,"withdrawn":42}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &MockBalanceService{
				GetUserBalanceFunc: tt.mockFunc,
			}

			handler := NewUserHandler(mockSvc)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest("GET", "/api/user/balance", nil)

			c.Set("userID", tt.userID)

			handler.GetBalance(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestWithdraw(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		body           string
		mockWithdraw   func(ctx context.Context, userID, order string, sum float64) (float64, error)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "unauthorized",
			userID:         "",
			body:           `{}`,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid json",
			userID:         "user-1",
			body:           `{invalid}`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:   "invalid order",
			userID: "user-1",
			body:   `{"order":"abc","sum":100}`,
			mockWithdraw: func(ctx context.Context, userID, order string, sum float64) (float64, error) {
				return 0, service.ErrInvalidOrder
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:   "insufficient funds",
			userID: "user-1",
			body:   `{"order":"123456","sum":1000}`,
			mockWithdraw: func(ctx context.Context, userID, order string, sum float64) (float64, error) {
				return 0, service.ErrInsufficientFunds
			},
			expectedStatus: http.StatusPaymentRequired, // 402
		},
		{
			name:   "internal error",
			userID: "user-1",
			body:   `{"order":"123456","sum":100}`,
			mockWithdraw: func(ctx context.Context, userID, order string, sum float64) (float64, error) {
				return 0, errors.New("db error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "success",
			userID: "user-1",
			body:   `{"order":"123456","sum":100}`,
			mockWithdraw: func(ctx context.Context, userID, order string, sum float64) (float64, error) {
				return 400, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"current":400}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &MockBalanceService{
				WithdrawFunc: tt.mockWithdraw,
			}

			h := &UserHandler{
				BalanceService: mockSvc,
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			if tt.userID != "" {
				c.Set("userID", tt.userID)
			}

			h.Withdraw(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestGetWithdrawals(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		userID       string
		mockFunc     func(ctx context.Context, userID string) ([]models.WithdrawalResponse, error)
		expectedCode int
		expectedBody string
	}{
		{
			name:   "unauthorized",
			userID: "",
			mockFunc: func(ctx context.Context, userID string) ([]models.WithdrawalResponse, error) {
				return nil, nil
			},
			expectedCode: http.StatusUnauthorized,
			expectedBody: "",
		},
		{
			name:   "no withdrawals",
			userID: "123",
			mockFunc: func(ctx context.Context, userID string) ([]models.WithdrawalResponse, error) {
				return []models.WithdrawalResponse{}, nil
			},
			expectedCode: http.StatusNoContent,
			expectedBody: "",
		},
		{
			name:   "internal error",
			userID: "123",
			mockFunc: func(ctx context.Context, userID string) ([]models.WithdrawalResponse, error) {
				return nil, errors.New("db error")
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "",
		},
		{
			name:   "success with withdrawals",
			userID: "123",
			mockFunc: func(ctx context.Context, userID string) ([]models.WithdrawalResponse, error) {
				return []models.WithdrawalResponse{
					{
						Order:       "2377225624",
						Sum:         500,
						ProcessedAt: time.Date(2020, 12, 9, 16, 9, 57, 0, time.FixedZone("MSK", 3*3600)).Format(time.RFC3339),
					},
				}, nil
			},
			expectedCode: http.StatusOK,
			expectedBody: `[{"order":"2377225624","sum":500,"processedAt":"2020-12-09T16:09:57+03:00"}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &MockBalanceService{
				GetWithdrawalsFunc: tt.mockFunc,
			}
			handler := &UserHandler{
				BalanceService: mockSvc,
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)

			if tt.userID != "" {
				c.Set("userID", tt.userID)
			}

			handler.GetWithdrawals(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}
