package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Guldana11/gophermart/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockBalanceService struct {
	GetUserBalanceFunc func(ctx context.Context, userID string) (float64, float64, error)
	WithdrawFunc       func(ctx context.Context, userID, order string, sum float64) (float64, error)
}

func (m *MockBalanceService) GetUserBalance(ctx context.Context, userID string) (float64, float64, error) {
	return m.GetUserBalanceFunc(ctx, userID)
}

func (m *MockBalanceService) Withdraw(ctx context.Context, userID, order string, sum float64) (float64, error) {
	return m.WithdrawFunc(ctx, userID, order, sum)
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

func TestUserHandler_Withdraw(t *testing.T) {
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
