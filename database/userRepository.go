package database

import (
	"context"

	"github.com/Guldana11/gophermart/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, login, password string) (*models.User, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
}
