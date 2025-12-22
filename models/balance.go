package models

import "time"

type BalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type WithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type WithdrawResponse struct {
	Current float64 `json:"current"`
}

type WithdrawalResponse struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processedAt"`
}
type Withdrawal struct {
	OrderID   string    `json:"order"`
	Sum       float64   `json:"sum"`
	CreatedAt time.Time `json:"createdAt"`
}
