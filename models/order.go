package models

import "time"

type Order struct {
	ID         int       `json:"id"`
	Number     string    `json:"number"`
	UserID     string    `json:"userId"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploadedAt"`
}

type OrderAccrualResponse struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual *int64 `json:"accrual,omitempty"`
}
