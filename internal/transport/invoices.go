package transport

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/AlexCorn999/transaction-system/internal/domain"
	"github.com/AlexCorn999/transaction-system/internal/logger"
)

func (s *APIServer) Invoice(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logger.LogError("invoice", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var invoice domain.Invoice
	if err := json.Unmarshal(data, &invoice); err != nil {
		logger.LogError("invoice", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.invoices.AddInvoice(r.Context(), &invoice); err != nil {
		switch {
		case errors.Is(err, domain.ErrIncorrectCurrency):
			logger.LogError("invoice", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		case errors.Is(err, domain.ErrIncorrectAmount):
			logger.LogError("invoice", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		case errors.Is(err, domain.ErrIncorrectWalletNumber):
			logger.LogError("invoice", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		default:
			logger.LogError("invoice", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
}
