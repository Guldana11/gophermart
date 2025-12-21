package service

import (
	"context"
	"regexp"

	"github.com/Guldana11/gophermart/database"
)

type BalanceService struct {
	repo database.UserRepository
}

func NewBalanceService(repo database.UserRepository) *BalanceService {
	return &BalanceService{repo: repo}
}

func (s *BalanceService) GetUserBalance(ctx context.Context, userID string) (current float64, withdrawn float64, err error) {
	return s.repo.GetUserPoints(ctx, userID)
}

var orderRegexp = regexp.MustCompile(`^\d+$`)

func (s *BalanceService) Withdraw(ctx context.Context, userID string, order string, sum float64) (float64, error) {

	if order == "" || !orderRegexp.MatchString(order) {
		return 0, ErrInvalidOrder
	}

	if sum <= 0 {
		return 0, ErrInvalidOrder
	}

	newBalance, err := s.repo.WithdrawPoints(ctx, userID, order, sum)
	if err != nil {
		if err == ErrInsufficientFunds {
			return 0, ErrInsufficientFunds
		}
		return 0, err
	}

	return newBalance, nil
}

type BalanceServiceType interface {
	GetUserBalance(ctx context.Context, userID string) (float64, float64, error)
	Withdraw(ctx context.Context, userID, order string, sum float64) (float64, error)
}
