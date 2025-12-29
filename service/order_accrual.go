package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Guldana11/gophermart/models"
)

var (
	ErrTooManyReq = errors.New("too many requests")
)

type LoyaltyService interface {
	GetOrderAccrual(ctx context.Context, orderNumber string) (*models.OrderAccrualResponse, error)
}

type loyaltyService struct {
	baseURL string
	client  *http.Client
}

func NewLoyaltyService(baseURL string) LoyaltyService {
	return &loyaltyService{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *loyaltyService) GetOrderAccrual(
	ctx context.Context,
	orderNumber string,
) (*models.OrderAccrualResponse, error) {

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/api/orders/%s", s.baseURL, orderNumber),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {

	case http.StatusOK:
		var res models.OrderAccrualResponse
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return nil, err
		}
		return &res, nil

	case http.StatusNoContent:
		return &models.OrderAccrualResponse{
			Order:  orderNumber,
			Status: "PROCESSING",
		}, nil

	case http.StatusNotFound:
		return &models.OrderAccrualResponse{
			Order:  orderNumber,
			Status: "INVALID",
		}, nil

	case http.StatusTooManyRequests:
		return nil, ErrTooManyReq

	default:
		return nil, errors.New("unexpected accrual service response")
	}
}
