package database

import (
	"context"
	"errors"
	"time"

	"github.com/Guldana11/gophermart/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, login, password string) (*models.User, error) {
	var exists bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE login=$1)", login).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("login already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	id := uuid.New().String()
	createdAt := time.Now()
	_, err = r.db.Exec(ctx,
		"INSERT INTO users (id, login, password_hash, created_at) VALUES ($1,$2,$3,$4)",
		id, login, string(hash), createdAt,
	)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:           id,
		Login:        login,
		PasswordHash: string(hash),
		CreatedAt:    createdAt,
	}, nil
}

func (r *UserRepo) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(ctx,
		"SELECT id, login, password_hash, created_at FROM users WHERE login=$1", login).
		Scan(&u.ID, &u.Login, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
