package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Guldana11/gophermart/database"
	"github.com/Guldana11/gophermart/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestUploadOrderHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	origCheck := database.CheckOrderExists
	origCreate := database.CreateOrder

	t.Cleanup(func() {
		database.CheckOrderExists = origCheck
		database.CreateOrder = origCreate
	})

	tests := []struct {
		name           string
		body           string
		userID         *int
		mockCheck      func() (string, bool, error)
		mockCreateErr  error
		expectedStatus int
	}{
		{
			name:           "unauthorized",
			body:           "79927398713",
			userID:         nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "empty body",
			body:           "",
			userID:         intPtr(1),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "non digit order",
			body:           "12ab34",
			userID:         intPtr(1),
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "invalid luhn",
			body:           "1234567890",
			userID:         intPtr(1),
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:   "order exists same user",
			body:   "79927398713",
			userID: intPtr(1),
			mockCheck: func() (string, bool, error) {
				return "", true, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "order exists other user",
			body:   "79927398713",
			userID: intPtr(1),
			mockCheck: func() (string, bool, error) {
				return "", true, nil
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:   "new order",
			body:   "79927398713",
			userID: intPtr(1),
			mockCheck: func() (string, bool, error) {
				return "", false, nil
			},
			expectedStatus: http.StatusAccepted,
		},
		{
			name:   "db error",
			body:   "79927398713",
			userID: intPtr(1),
			mockCheck: func() (string, bool, error) {
				return "", false, errStub{}
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			database.CheckOrderExists = func(_ context.Context, _ string) (string, bool, error) {
				if tt.mockCheck != nil {
					return tt.mockCheck()
				}
				return "", false, nil
			}

			database.CreateOrder = func(_ context.Context, _ models.Order) error {
				return tt.mockCreateErr
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req := httptest.NewRequest(
				http.MethodPost,
				"/api/user/orders",
				strings.NewReader(tt.body),
			)
			c.Request = req

			if tt.userID != nil {
				c.Set("userID", *tt.userID)
			}

			UploadOrderHandler(c)

			require.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func intPtr(i int) *int {
	return &i
}

type errStub struct{}

func (errStub) Error() string {
	return "db error"
}

func TestGetOrdersHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	origGetOrders := database.GetOrdersByUser
	t.Cleanup(func() {
		database.GetOrdersByUser = origGetOrders
	})

	userID := "9da96116-0301-478b-9583-f0335368d51a"
	sampleOrders := []models.Order{
		{
			Number:     "1234567890",
			Status:     "PROCESSED",
			Accrual:    100,
			UploadedAt: time.Now(),
		},
		{
			Number:     "9876543210",
			Status:     "PROCESSING",
			UploadedAt: time.Now().Add(-time.Hour),
		},
	}

	tests := []struct {
		name           string
		userID         *string
		mockOrders     func() ([]models.Order, error)
		expectedStatus int
	}{
		{
			name:           "unauthorized",
			userID:         nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:   "no orders",
			userID: &userID,
			mockOrders: func() ([]models.Order, error) {
				return []models.Order{}, nil
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:   "orders exist",
			userID: &userID,
			mockOrders: func() ([]models.Order, error) {
				return sampleOrders, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "db error",
			userID: &userID,
			mockOrders: func() ([]models.Order, error) {
				return nil, errStub{}
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// мокируем БД
			database.GetOrdersByUser = func(_ context.Context, _ string) ([]models.Order, error) {
				if tt.mockOrders != nil {
					return tt.mockOrders()
				}
				return nil, nil
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
			c.Request = req

			if tt.userID != nil {
				c.Set("userID", *tt.userID)
			}

			GetOrdersHandler(c)

			require.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
