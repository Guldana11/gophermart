package service

import (
	"context"
	"regexp"

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

var orderRegexp = regexp.MustCompile(`^\d{1,20}$`)

func (s *BalanceService) Withdraw(ctx context.Context, userID string, order string, sum float64) (float64, error) {
	if order == "" || !orderRegexp.MatchString(order) {
		return 0, ErrInvalidOrder
	}
	if sum <= 0 {
		return 0, ErrInvalidOrder
	}

	newBalance, err := s.repo.WithdrawPoints(ctx, userID, order, sum)
	if err != nil {
		return 0, err
	}

	return newBalance, nil
}

func (s *BalanceService) GetWithdrawals(ctx context.Context, userID string) ([]models.Withdrawal, error) {
	rows, err := s.repo.GetUserWithdrawals(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := make([]models.Withdrawal, len(rows))
	for i, w := range rows {
		res[i] = models.Withdrawal{
			Order:       w.Order,
			Sum:         w.Sum,
			ProcessedAt: w.ProcessedAt,
		}
	}
	return res, nil
}

func (s *BalanceService) SaveWithdrawal(ctx context.Context, userID string, order string, sum float64) error {
	return s.repo.SaveWithdrawal(ctx, userID, order, sum)
}

type BalanceServiceType interface {
	GetUserBalance(ctx context.Context, userID string) (float64, float64, error)
	Withdraw(ctx context.Context, userID, order string, sum float64) (float64, error)
	GetWithdrawals(ctx context.Context, userID string) ([]models.Withdrawal, error)
	SaveWithdrawal(ctx context.Context, userID, order string, sum float64) error
}
