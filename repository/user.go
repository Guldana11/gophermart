package repository

import (
	"context"

	"github.com/Guldana11/gophermart/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, login, password string) (*models.User, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	GetUserPoints(ctx context.Context, userID string) (float64, float64, error)
	WithdrawPoints(ctx context.Context, userID, order string, sum float64) (float64, error)
	GetUserWithdrawals(ctx context.Context, userID string) ([]models.Withdrawal, error)
}
