package transport

import (
	"context"

	"github.com/AlexCorn999/transaction-system/internal/logger"
)

// orderProcessing moves the order to a new status every 2 seconds.
func (s *APIServer) orderProcess() {

	// get order numbers up to 2 pcs from the system if their status is not SUCCESS
	invoices, err := s.orderProcessing.GetInvoiceStatus(context.Background())
	if err != nil {
		logger.LogError("orderProcess", err)
		return
	}

	for _, invoice := range invoices {
		if err := s.orderProcessing.UpdateInvoice(context.Background(), invoice); err != nil {
			logger.LogError("orderProcess", err)
			return
		}
	}
}
