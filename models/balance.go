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

type Withdrawal struct {
	UserId      string    `json:"user_id"`
	OrderNumber string    `json:"orderNumber"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processedAt"`
}
