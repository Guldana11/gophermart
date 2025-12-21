package database

import (
	"context"
	"errors"
	"time"

	"github.com/Guldana11/gophermart/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var ErrInsufficientFunds = errors.New("insufficient funds")
var ErrInvalidOrder = errors.New("invalid order")

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

func (r *UserRepo) GetUserPoints(ctx context.Context, userID string) (float64, float64, error) {
	var current, withdrawn float64
	err := r.db.QueryRow(ctx,
		`SELECT current_balance, withdrawn_points 
		 FROM user_points 
		 WHERE user_id = $1`, userID).
		Scan(&current, &withdrawn)
	if err != nil {
		return 0, 0, err
	}
	return current, withdrawn, nil
}

func (r *UserRepo) WithdrawPoints(
	ctx context.Context,
	userID string,
	order string,
	sum float64,
) (float64, error) {

	var current float64

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var exists bool
	err = tx.QueryRow(ctx,
		`SELECT EXISTS (
			SELECT 1 FROM withdrawals WHERE order_id = $1
		)`,
		order,
	).Scan(&exists)

	if err != nil {
		return 0, err
	}
	if exists {
		return 0, ErrInvalidOrder
	}

	err = tx.QueryRow(ctx,
		`SELECT current_balance
		 FROM user_points
		 WHERE user_id = $1
		 FOR UPDATE`,
		userID,
	).Scan(&current)

	if err != nil {
		return 0, err
	}

	if sum > current {
		return current, ErrInsufficientFunds
	}

	_, err = tx.Exec(ctx,
		`UPDATE user_points
		 SET current_balance = current_balance - $1,
		     withdrawn_points = withdrawn_points + $1
		 WHERE user_id = $2`,
		sum, userID,
	)
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO withdrawals (order_id, user_id, sum)
		 VALUES ($1, $2, $3)`,
		order, userID, sum,
	)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return current - sum, nil
}
