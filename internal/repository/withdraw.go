package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/AlexCorn999/transaction-system/internal/domain"
)

// Withdraw добавляет списания бонусов пользователя.
func (s *Storage) Withdraw(tx *sql.Tx, withdraw domain.Withdraw) error {
	_, err := tx.Exec("INSERT INTO withdrawals (wallet_number, currency, uploaded_at, amount) values ($1, $2, $3, $4)",
		withdraw.WalletNumber, withdraw.Currency, withdraw.UploadedAt, withdraw.Amount)
	if err != nil {
		return fmt.Errorf("postgreSQL: withdraw %s", err)
	}

	return nil
}

// Balance возвращает весь баланс кошелька.
func (s *Storage) Balance(ctx context.Context, withdraw *domain.Withdraw) (float32, error) {
	var nullableBalance sql.NullFloat64
	err := s.DB.QueryRowContext(ctx, "SELECT SUM(amount) FROM invoices WHERE wallet_number=$1 AND currency=$2", withdraw.WalletNumber, withdraw.Currency).
		Scan(&nullableBalance)
	if err != nil {
		return 0, fmt.Errorf("postgreSQL: balance %s", err)
	}
	if !nullableBalance.Valid {
		return 0, nil
	}

	balance := float32(nullableBalance.Float64)
	return balance, nil
}

// WithdrawBalance возвращает сумму списанных денег.
func (s *Storage) WithdrawBalance(ctx context.Context, withdraw *domain.Withdraw) (float32, error) {
	var nullableBalance sql.NullFloat64
	err := s.DB.QueryRowContext(ctx, "SELECT SUM(amount) FROM withdrawals WHERE wallet_number=$1 AND currency=$2", withdraw.WalletNumber, withdraw.Currency).
		Scan(&nullableBalance)
	if err != nil {
		return 0, fmt.Errorf("postgreSQL: withdrawBalance %s", err)
	}
	if !nullableBalance.Valid {
		return 0, nil
	}

	balance := float32(nullableBalance.Float64)
	return balance, nil
}
