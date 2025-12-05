package service

import (
	"context"

	"github.com/Guldana11/gophermart/models"
)

type UserServiceInterface interface {
	Register(ctx context.Context, login, password string) (*models.User, error)
	Login(ctx context.Context, login, password string) (*models.User, error)
}
