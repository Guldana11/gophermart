package repository

import (
	"context"

	"github.com/Guldana11/gophermart/models"
)

type OrderRepository interface {
	CheckOrderExists(ctx context.Context, orderNumber string) (string, bool, error)
	CreateOrder(ctx context.Context, order models.Order) error
	GetOrdersByUser(ctx context.Context, userID string) ([]models.Order, error)
}
