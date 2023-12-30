package transport

import (
	"net/http"
)

func (s *APIServer) Invoice(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hey"))
}
