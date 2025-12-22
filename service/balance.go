package service

import (
	"context"
	"regexp"
	"time"

	"github.com/Guldana11/gophermart/models"
	"github.com/Guldana11/gophermart/repository"
)

type BalanceService struct {
	repo repository.UserRepository
}

func NewBalanceService(repo repository.UserRepository) *BalanceService {
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

func (s *BalanceService) GetWithdrawals(ctx context.Context, userID string) ([]models.WithdrawalResponse, error) {
	rows, err := s.repo.GetUserWithdrawals(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := make([]models.WithdrawalResponse, len(rows))
	for i, w := range rows {
		res[i] = models.WithdrawalResponse{
			Order:       w.OrderID,
			Sum:         w.Sum,
			ProcessedAt: w.CreatedAt.Format(time.RFC3339),
		}
	}
	return res, nil
}

type BalanceServiceType interface {
	GetUserBalance(ctx context.Context, userID string) (float64, float64, error)
	Withdraw(ctx context.Context, userID, order string, sum float64) (float64, error)
	GetWithdrawals(ctx context.Context, userID string) ([]models.WithdrawalResponse, error)
}
