package transport

import (
	"net/http"
)

func (s *APIServer) BalanceFrozen(w http.ResponseWriter, r *http.Request) {
	// balance, err := s.withdraw.Balance(r.Context())
	// if err != nil {
	// 	logger.LogError("balance", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// balanceJSON, err := json.Marshal(balance)
	// if err != nil {
	// 	logger.LogError("balance", err)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//w.Write(balanceJSON)
}
