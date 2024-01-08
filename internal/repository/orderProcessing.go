package repository

import (
	"context"
	"fmt"

	"github.com/AlexCorn999/transaction-system/internal/domain"
)

// GetInvoiceStatus receives all invoices with CREATED statuses.
func (s *Storage) GetInvoiceStatus(ctx context.Context) ([]domain.Invoice, error) {
	var invoices []domain.Invoice
	rows, err := s.DB.QueryContext(ctx, "SELECT * FROM invoices WHERE status NOT IN ('SUCCESS') LIMIT 2")
	if err != nil {
		return nil, fmt.Errorf("postgreSQL: getInvoiceStatus %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var invoice domain.Invoice
		err := rows.Scan(&invoice.WalletNumber, &invoice.Currency, &invoice.UploadedAt, &invoice.Amount, &invoice.Status, &invoice.UserID)
		if err != nil {
			return nil, fmt.Errorf("postgreSQL: getInvoiceStatus %w", err)
		}
		invoices = append(invoices, invoice)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgreSQL: getInvoiceStatus %w", err)
	}

	return invoices, nil
}

// UpdateInvoice updates invoice status.
func (s *Storage) UpdateInvoice(ctx context.Context, invoice domain.Invoice) error {
	_, err := s.DB.ExecContext(ctx, "UPDATE invoices SET status=$1 WHERE wallet_number=$2 AND currency=$3 AND amount=$4 AND user_id=$5", domain.Success, invoice.WalletNumber, invoice.Currency, invoice.Amount, &invoice.UserID)
	if err != nil {
		return fmt.Errorf("postgreSQL: updateInvoice %w", err)
	}
	return nil
}
