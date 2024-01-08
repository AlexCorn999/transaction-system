package service

import (
	"context"

	"github.com/AlexCorn999/transaction-system/internal/domain"
)

type OrderProcessingRepository interface {
	GetInvoiceStatus(ctx context.Context) ([]domain.Invoice, error)
	UpdateInvoice(ctx context.Context, invoice domain.Invoice) error
}

type OrderProcessing struct {
	repo OrderProcessingRepository
}

func NewOrderProcessing(repo OrderProcessingRepository) *OrderProcessing {
	return &OrderProcessing{
		repo: repo,
	}
}

// GetInvoiceStatus receives all invoices with CREATED statuses.
func (o *OrderProcessing) GetInvoiceStatus(ctx context.Context) ([]domain.Invoice, error) {
	return o.repo.GetInvoiceStatus(ctx)
}

// UpdateInvoice updates invoice status.
func (o *OrderProcessing) UpdateInvoice(ctx context.Context, invoice domain.Invoice) error {
	return o.repo.UpdateInvoice(ctx, invoice)
}
