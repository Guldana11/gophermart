package service

import (
	"context"
	"errors"

	"github.com/Guldana11/gophermart/models"
)

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrTooManyReq    = errors.New("too many requests")
)

type LoyaltyServiceType interface {
	GetOrderAccrual(ctx context.Context, orderNumber string) (*models.OrderAccrualResponse, error)
}

type LoyaltyService struct {
}

func NewLoyaltyService() *LoyaltyService {
	return &LoyaltyService{}
}

func (s *LoyaltyService) GetOrderAccrual(ctx context.Context, orderNumber string) (*models.OrderAccrualResponse, error) {
	if orderNumber == "" {
		return nil, errors.New("order number empty")
	}

	return &models.OrderAccrualResponse{
		Order:   orderNumber,
		Status:  "REGISTERED",
		Accrual: 0,
	}, nil
}
