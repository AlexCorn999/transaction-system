package transport

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/AlexCorn999/transaction-system/internal/domain"
	"github.com/AlexCorn999/transaction-system/internal/logger"
)

func (s *APIServer) Withdraw(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logger.LogError("withdraw", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var withdraw domain.Withdraw
	if err := json.Unmarshal(data, &withdraw); err != nil {
		logger.LogError("withdraw", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.withdraw.Withdraw(r.Context(), withdraw); err != nil {
		switch {
		case errors.Is(err, domain.ErrIncorrectCurrency):
			logger.LogError("withdraw", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		case errors.Is(err, domain.ErrIncorrectAmount):
			logger.LogError("withdraw", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		case errors.Is(err, domain.ErrIncorrectAmount):
			logger.LogError("withdraw", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		case errors.Is(err, domain.ErrNoMoney):
			logger.LogError("withdraw", err)
			w.WriteHeader(http.StatusPaymentRequired)
			return
		default:
			logger.LogError("withdraw", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
}
