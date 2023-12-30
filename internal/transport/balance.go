package transport

import (
	"encoding/json"
	"net/http"

	"github.com/AlexCorn999/transaction-system/internal/logger"
)

func (s *APIServer) BalanceActual(w http.ResponseWriter, r *http.Request) {
	balance, err := s.withdraw.Balance(r.Context())
	if err != nil {
		logger.LogError("balance", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	balanceJSON, err := json.Marshal(balance)
	if err != nil {
		logger.LogError("balance", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(balanceJSON)
}

func (s *APIServer) BalanceHold(w http.ResponseWriter, r *http.Request) {
	balance, err := s.withdraw.Balance(r.Context())
	if err != nil {
		logger.LogError("balance", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	balanceJSON, err := json.Marshal(balance)
	if err != nil {
		logger.LogError("balance", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(balanceJSON)
}
