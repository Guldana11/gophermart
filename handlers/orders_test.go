package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Guldana11/gophermart/models"
	"github.com/Guldana11/gophermart/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestOrderHandler_GetOrdersHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type fields struct {
		orderService service.OrderService
	}

	tests := []struct {
		name           string
		fields         fields
		userID         string
		mockFunc       func(ctx context.Context, userID string) ([]models.Order, error)
		expectedStatus int
		expectedBody   []models.Order
	}{
		{
			name:   "500 internal server error",
			userID: "user1",
			mockFunc: func(ctx context.Context, userID string) ([]models.Order, error) {
				return nil, errors.New("unknown error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
		{
			name:   "200 empty orders",
			userID: "user1",
			mockFunc: func(ctx context.Context, userID string) ([]models.Order, error) {
				return nil, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   []models.Order{},
		},
		{
			name:   "200 orders returned",
			userID: "user1",
			mockFunc: func(ctx context.Context, userID string) ([]models.Order, error) {
				return []models.Order{
					{UserID: userID, Number: "12345678903"},
					{UserID: userID, Number: "79927398713"},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody: []models.Order{
				{UserID: "user1", Number: "12345678903"},
				{UserID: "user1", Number: "79927398713"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/api/user/orders", nil)
			if tt.userID != "" {
				c.Set("userID", tt.userID)
			}

			mockService := &MockOrderService{
				GetOrdersFunc: tt.mockFunc,
			}

			h := &OrderHandler{
				orderService: mockService,
			}

			h.GetOrdersHandler(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var got []models.Order
				err := json.Unmarshal(w.Body.Bytes(), &got)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, got)
			}
		})
	}
}

func TestOrderHandler_UploadOrderHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		body           string
		mockFunc       func(ctx context.Context, userID, orderNumber string) error
		expectedStatus int
	}{
		{
			name:   "401 unauthorized if no userID",
			userID: "",
			body:   "12345678903",
			mockFunc: func(ctx context.Context, userID, orderNumber string) error {
				return nil
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:   "400 bad request if body empty",
			userID: "user1",
			body:   "",
			mockFunc: func(ctx context.Context, userID, orderNumber string) error {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "422 unprocessable entity if number invalid (non-numeric)",
			userID: "user1",
			body:   "abc123",
			mockFunc: func(ctx context.Context, userID, orderNumber string) error {
				return nil
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:   "422 unprocessable entity if number invalid (fails Luhn)",
			userID: "user1",
			body:   "12345678901",
			mockFunc: func(ctx context.Context, userID, orderNumber string) error {
				return nil
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:   "202 accepted new number",
			userID: "user1",
			body:   "79927398713",
			mockFunc: func(ctx context.Context, userID, orderNumber string) error {
				return nil
			},
			expectedStatus: http.StatusAccepted,
		},
		{
			name:   "200 already uploaded by self",
			userID: "user1",
			body:   "79927398713",
			mockFunc: func(ctx context.Context, userID, orderNumber string) error {
				return service.ErrAlreadyUploadedSelf
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "409 already uploaded by other",
			userID: "user1",
			body:   "79927398713",
			mockFunc: func(ctx context.Context, userID, orderNumber string) error {
				return service.ErrAlreadyUploadedOther
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:   "500 internal server error",
			userID: "user1",
			body:   "79927398713",
			mockFunc: func(ctx context.Context, userID, orderNumber string) error {
				return errors.New("unknown error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest("POST", "/api/user/orders", bytes.NewBufferString(tt.body))
			if tt.userID != "" {
				c.Set("userID", tt.userID)
			}

			mockService := &MockOrderService{
				UploadFunc: tt.mockFunc,
			}

			h := &OrderHandler{
				orderService: mockService,
			}

			h.UploadOrderHandler(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

type MockOrderService struct {
	UploadFunc    func(ctx context.Context, userID, orderNumber string) error
	GetOrdersFunc func(ctx context.Context, userID string) ([]models.Order, error)
}

func (m *MockOrderService) UploadOrder(ctx context.Context, userID, orderNumber string) error {
	return m.UploadFunc(ctx, userID, orderNumber)
}

func (m *MockOrderService) GetOrders(ctx context.Context, userID string) ([]models.Order, error) {
	if m.GetOrdersFunc != nil {
		return m.GetOrdersFunc(ctx, userID)
	}
	return nil, nil
}
