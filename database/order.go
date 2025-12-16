package database

import (
	"context"
	"errors"

	"github.com/Guldana11/gophermart/models"
	"github.com/jackc/pgx/v5"
)

func checkOrderExists(ctx context.Context, orderNumber string) (string, bool, error) {
	var userID string

	err := DB.QueryRow(
		ctx,
		"SELECT user_id FROM orders WHERE number = $1",
		orderNumber,
	).Scan(&userID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", false, nil
		}
		return "", false, err
	}

	return userID, true, nil
}

func createOrder(ctx context.Context, order models.Order) error {
	_, err := DB.Exec(
		ctx,
		"INSERT INTO orders (user_id, number) VALUES ($1, $2)",
		order.UserID,
		order.Number,
	)
	return err
}

var CheckOrderExists = checkOrderExists
var CreateOrder = createOrder

var GetOrdersByUser = func(ctx context.Context, userID string) ([]models.Order, error) {
	rows, err := DB.Query(ctx,
		`SELECT number, status, accrual, uploaded_at 
         FROM orders 
         WHERE user_id=$1 
         ORDER BY uploaded_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.Number, &o.Status, &o.Accrual, &o.UploadedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	if len(orders) == 0 {
		return nil, nil
	}
	return orders, nil
}
