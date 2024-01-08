package transport

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/AlexCorn999/transaction-system/internal/domain"
	"github.com/AlexCorn999/transaction-system/internal/logger"
	"github.com/AlexCorn999/transaction-system/internal/repository"
)

// SighUp is responsible for user registration.
func (s *APIServer) SighUp(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logger.LogError("signUp", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var usr domain.SighUpAndInInput
	if err := json.Unmarshal(data, &usr); err != nil {
		logger.LogError("signUp", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := usr.Validate(); err != nil {
		logger.LogError("signUp", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.users.SignUp(r.Context(), usr); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			w.WriteHeader(http.StatusConflict)
			return
		} else {
			logger.LogError("signUp", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}

	// atomatic user authorization
	r.Body = io.NopCloser(bytes.NewBuffer(data))
	s.SighIn(w, r)
}

// SighIn is responsible for authorizing the user. Issues an access token.
func (s *APIServer) SighIn(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logger.LogError("signIn", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var usr domain.SighUpAndInInput
	if err := json.Unmarshal(data, &usr); err != nil {
		logger.LogError("signIn", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := usr.Validate(); err != nil {
		logger.LogError("signIn", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := s.users.SignIn(r.Context(), usr)
	if err != nil {
		// user not found
		if errors.Is(err, domain.ErrUserNotFound) {
			logger.LogError("signIn", err)
			w.WriteHeader(http.StatusUnauthorized)
		}
		logger.LogError("signIn", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusOK)
}
