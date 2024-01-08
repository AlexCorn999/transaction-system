package transport

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/AlexCorn999/transaction-system/internal/domain"
	"github.com/AlexCorn999/transaction-system/internal/logger"
)

// Invoice credits money to the user's account.
func (s *APIServer) Invoice(data []byte, userID string) {
	var invoice domain.Invoice
	if err := json.Unmarshal(data, &invoice); err != nil {
		logger.LogError("invoice", err)
		return
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		logger.LogError("invoice", err)
		return
	}

	ctx := context.WithValue(context.Background(), domain.UserIDKeyForContext, int64(id))

	if err := s.money.Invoice(ctx, &invoice); err != nil {
		switch {
		case errors.Is(err, domain.ErrIncorrectCurrency):
			logger.LogError("invoice", err)
			//w.WriteHeader(http.StatusUnprocessableEntity)
			return
		case errors.Is(err, domain.ErrIncorrectAmount):
			logger.LogError("invoice", err)
			//w.WriteHeader(http.StatusUnprocessableEntity)
			return
		case errors.Is(err, domain.ErrIncorrectWalletNumber):
			logger.LogError("invoice", err)
			//w.WriteHeader(http.StatusUnprocessableEntity)
			return
		default:
			logger.LogError("invoice", err)
			//w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	//w.WriteHeader(http.StatusAccepted)
}

// Withdraw withdraws money from the user's account.
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

	if err := s.money.Withdraw(r.Context(), withdraw); err != nil {
		switch {
		case errors.Is(err, domain.ErrIncorrectCurrency):
			logger.LogError("withdraw", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		case errors.Is(err, domain.ErrIncorrectAmount):
			logger.LogError("withdraw", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		case errors.Is(err, domain.ErrIncorrectWalletNumber):
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

// BalanceActual returns the user's wallet balance with success status.
func (s *APIServer) BalanceActual(w http.ResponseWriter, r *http.Request) {
	balance, err := s.money.Balance(r.Context())
	if err != nil {
		logger.LogError("balanceActual", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	balanceJSON, err := json.Marshal(balance)
	if err != nil {
		logger.LogError("balanceActual", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(balanceJSON)
}

// BalanceFrozen displays the user's balance in the created status.
func (s *APIServer) BalanceFrozen(w http.ResponseWriter, r *http.Request) {
	balance, err := s.money.BalanceFrozen(r.Context())
	if err != nil {
		logger.LogError("balanceFrozen", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	balanceJSON, err := json.Marshal(balance)
	if err != nil {
		logger.LogError("balanceFrozen", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(balanceJSON)
}
