package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/AlexCorn999/transaction-system/internal/domain"
)

// Invoice credits money to the user's account.
func (s *Storage) Invoice(ctx context.Context, invoice *domain.Invoice) error {
	row := s.DB.QueryRowContext(ctx, "SELECT currency FROM invoices WHERE NOT user_id=$1 AND wallet_number=$2", invoice.UserID, invoice.WalletNumber)

	var currency string

	if err := row.Scan(&currency); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("postgreSQL: invoice %w", err)
		}
	}

	// if the wallet number already exists, but for another user
	if len(currency) > 0 {
		return domain.ErrIncorrectWalletNumber
	}

	_, err := s.DB.ExecContext(ctx, "INSERT INTO invoices (wallet_number, currency, uploaded_at, amount, status, user_id) values ($1, $2, $3, $4, $5, $6)",
		invoice.WalletNumber, invoice.Currency, invoice.UploadedAt, invoice.Amount, invoice.Status, invoice.UserID)
	if err != nil {
		return fmt.Errorf("postgreSQL: invoice %w", err)
	}

	return nil
}

// Withdraw debits money from the user's wallet.
func (s *Storage) Withdraw(ctx context.Context, withdraw domain.Withdraw) error {
	// currency account verification
	tx := extractTx(ctx)
	var walletNumber string
	row := tx.QueryRow("SELECT wallet_number FROM invoices WHERE user_id=$1 AND currency=$2", withdraw.UserID, withdraw.Currency)
	if err := row.Scan(&walletNumber); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrIncorrectWalletNumber
		} else {
			return fmt.Errorf("postgreSQL: invoice %w", err)
		}
	}

	_, err := tx.Exec("INSERT INTO withdrawals (wallet_number, currency, uploaded_at, amount, user_id) values ($1, $2, $3, $4, $5)",
		walletNumber, withdraw.Currency, withdraw.UploadedAt, withdraw.Amount, withdraw.UserID)
	if err != nil {
		return fmt.Errorf("postgreSQL: withdraw %w", err)
	}

	return nil
}

// CheckWallet checks if the user has an account and returns the user's account id.
func (s *Storage) CheckWallet(ctx context.Context, withdraw domain.Withdraw) (int, error) {
	tx := extractTx(ctx)
	var userID int
	row := tx.QueryRow("SELECT user_id FROM invoices WHERE wallet_number=$1 AND currency=$2 AND NOT user_id=$3", withdraw.WalletNumber, withdraw.Currency, withdraw.UserID)
	if err := row.Scan(&userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, domain.ErrIncorrectWalletNumber
		} else {
			return 0, fmt.Errorf("postgreSQL: invoice %w", err)
		}
	}
	return userID, nil
}

// Invoice credits money to another user's account.
func (s *Storage) InvoiceToUser(ctx context.Context, invoice *domain.Invoice) error {
	tx := extractTx(ctx)
	_, err := tx.Exec("INSERT INTO invoices (wallet_number, currency, uploaded_at, amount, status, user_id) values ($1, $2, $3, $4, $5, $6)",
		invoice.WalletNumber, invoice.Currency, invoice.UploadedAt, invoice.Amount, invoice.Status, invoice.UserID)
	if err != nil {
		return fmt.Errorf("postgreSQL: invoice %w", err)
	}

	return nil
}

// Balance returns the user's wallet balance with success status.
func (s *Storage) Balance(ctx context.Context, withdraw *domain.Withdraw) (float32, error) {
	tx := extractTx(ctx)
	var nullableBalance sql.NullFloat64
	// поменять статус на success TODO
	err := tx.QueryRow("SELECT SUM(amount) FROM invoices WHERE currency=$1 AND user_id=$2 AND status=$3", withdraw.Currency, withdraw.UserID, domain.Created).
		Scan(&nullableBalance)
	if err != nil {
		return 0, fmt.Errorf("postgreSQL: balance %w", err)
	}
	if !nullableBalance.Valid {
		return 0, nil
	}

	balance := float32(nullableBalance.Float64)
	return balance, nil
}

// WithdrawBalance returns the amount of debited money of the user.
func (s *Storage) WithdrawBalance(ctx context.Context, withdraw *domain.Withdraw) (float32, error) {
	tx := extractTx(ctx)
	var nullableBalance sql.NullFloat64
	err := tx.QueryRow("SELECT SUM(amount) FROM withdrawals WHERE currency=$1 AND user_id=$2", withdraw.Currency, withdraw.UserID).
		Scan(&nullableBalance)
	if err != nil {
		return 0, fmt.Errorf("postgreSQL: withdrawBalance %w", err)
	}
	if !nullableBalance.Valid {
		return 0, nil
	}

	balance := float32(nullableBalance.Float64)
	return balance, nil
}

// Balance returns the user's wallet balance with success status.
func (s *Storage) BalanceActual(userID int64) ([]domain.BalanceOutput, error) {
	// поменять статус на success TODO
	rows, err := s.DB.Query("SELECT i.currency, SUM(i.amount) - COALESCE(w.total_amount, 0) AS difference FROM invoices AS i LEFT JOIN (SELECT currency, SUM(amount) AS total_amount FROM withdrawals WHERE user_id=$1 GROUP BY currency) AS w ON i.currency = w.currency WHERE i.user_id=$2 AND i.status=$3 GROUP BY i.currency, w.total_amount", userID, userID, domain.Created)
	if err != nil {
		return nil, fmt.Errorf("postgreSQL: balance %w", err)
	}
	defer rows.Close()

	balanceOutput := make([]domain.BalanceOutput, 0)

	for rows.Next() {
		balance := domain.BalanceOutput{}
		if err := rows.Scan(&balance.Currency, &balance.Amount); err != nil {
			return nil, err
		}
		balanceOutput = append(balanceOutput, balance)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return balanceOutput, nil
}
