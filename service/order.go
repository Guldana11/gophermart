package service

import (
	"context"
	"strings"

	"github.com/Guldana11/gophermart/models"
	"github.com/Guldana11/gophermart/repository"
)

type OrderService interface {
	UploadOrder(ctx context.Context, userID, orderNumber string) error
	GetOrders(ctx context.Context, userID string) ([]models.Order, error)
}

type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) UploadOrder(ctx context.Context, userID, orderNumber string) error {
	orderNumber = strings.TrimSpace(orderNumber)

	if len(orderNumber) <= 1 {
		return ErrInvalidOrder
	}

	for _, r := range orderNumber {
		if r < '0' || r > '9' {
			return ErrInvalidOrder
		}
	}

	if !CheckLuhn(orderNumber) {
		return ErrInvalidOrder
	}

	existingUserID, exists, err := s.repo.CheckOrderExists(ctx, orderNumber)
	if err != nil {
		return err
	}

	if exists {
		if existingUserID == userID {
			return ErrAlreadyUploadedSelf
		}
		return ErrAlreadyUploadedOther
	}

	order := models.Order{
		UserID: userID,
		Number: orderNumber,
	}

	return s.repo.CreateOrder(ctx, order)
}

func (s *orderService) GetOrders(ctx context.Context, userID string) ([]models.Order, error) {
	return s.repo.GetOrdersByUser(ctx, userID)
}

func CheckLuhn(number string) bool {
	sum := 0
	double := false

	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')
		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		double = !double
	}

	return sum%10 == 0
}
