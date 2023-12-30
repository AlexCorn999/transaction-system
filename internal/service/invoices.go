package service

import (
	"context"
	"strings"
	"time"

	"github.com/AlexCorn999/transaction-system/internal/domain"
)

type InvoicesRepository interface {
	AddInvoice(ctx context.Context, invoice *domain.Invoice) error
}

type Invoices struct {
	repo InvoicesRepository
}

func NewInvoices(repo InvoicesRepository) *Invoices {
	return &Invoices{
		repo: repo,
	}
}

// AddInvoice зачисляет средства на счет кошелька.
func (o *Invoices) AddInvoice(ctx context.Context, invoice *domain.Invoice) error {
	if len(strings.TrimSpace(invoice.WalletNumber)) == 0 {
		return domain.ErrIncorrectWalletNumber
	}

	if invoice.Amount <= 0 {
		return domain.ErrIncorrectAmount
	}

	if _, ok := domain.Currency[invoice.Currency]; !ok {
		return domain.ErrIncorrectCurrency

	}

	invoice.UploadedAt = time.Now().Format(time.RFC3339)
	invoice.Status = domain.Created

	return o.repo.AddInvoice(ctx, invoice)
}
