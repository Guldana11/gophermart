package database

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/Guldana11/gophermart/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestUserRepo_Withdraw(t *testing.T) {
	ctx := context.Background()
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/loyalty?sslmode=disable"
	}

	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	r := &UserRepo{db: db}

	userID := "00000000-0000-0000-0000-000000000001"
	_, _ = db.Exec(ctx, "DELETE FROM user_points WHERE user_id=$1", userID)
	_, _ = db.Exec(ctx, "DELETE FROM withdrawals WHERE user_id=$1", userID)
	_, _ = db.Exec(ctx, `INSERT INTO user_points (user_id, current_balance, withdrawn_points) VALUES ($1, $2, $3)`, userID, 500, 0)

	tests := []struct {
		name        string
		order       string
		sum         float64
		wantErr     bool
		expectedErr error
	}{
		{
			name:  "success withdraw",
			order: "test-order-1",
			sum:   100,
		},
		{
			name:        "insufficient funds",
			order:       "test-order-2",
			sum:         1000,
			wantErr:     true,
			expectedErr: service.ErrInsufficientFunds,
		},
		{
			name:        "duplicate order",
			order:       "test-order-1",
			sum:         50,
			wantErr:     true,
			expectedErr: service.ErrInvalidOrder,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := r.Withdraw(ctx, userID, tt.order, tt.sum)
			if (err != nil) != tt.wantErr {
				t.Errorf("Withdraw() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && !errors.Is(err, tt.expectedErr) {
				t.Errorf("Withdraw() error = %v, expected %v", err, tt.expectedErr)
			}
		})
	}

	current, _, err := r.GetUserPoints(ctx, userID)
	if err != nil {
		t.Fatal(err)
	}
	if current != 400 {
		t.Errorf("unexpected balance: got %v, want %v", current, 400)
	}
}
