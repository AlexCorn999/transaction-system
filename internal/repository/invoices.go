package repository

import (
	"context"
	"fmt"

	"github.com/AlexCorn999/transaction-system/internal/domain"
)

// AddInvoice зачисляет средства на счет кошелька.
func (s *Storage) AddInvoice(ctx context.Context, invoice *domain.Invoice) error {
	_, err := s.DB.ExecContext(ctx, "INSERT INTO invoices (wallet_number, currency, uploaded_at, amount, status) values ($1, $2, $3, $4, $5)",
		invoice.WalletNumber, invoice.Currency, invoice.UploadedAt, invoice.Amount, invoice.Status)
	if err != nil {
		return fmt.Errorf("postgreSQL: addInvoice %s", err)
	}

	return nil
}
