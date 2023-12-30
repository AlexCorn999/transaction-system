package transport

import (
	"context"
	"net/http"

	"github.com/AlexCorn999/transaction-system/internal/domain"
	"github.com/AlexCorn999/transaction-system/internal/logger"
)

// authMiddleware выполняет функцию middleware авторизации.
// Получает токен из запроса и передает в контекст userID который совершает данный запрос.
func (s *APIServer) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getTokenFromRequest(r)
		if err != nil {
			logger.LogError("authMiddleware", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userID, err := s.users.ParseToken(r.Context(), token)
		if err != nil {
			logger.LogError("authMiddleware", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), domain.UserIDKeyForContext, userID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// getTokenFromRequest получает token из cookie.
func getTokenFromRequest(r *http.Request) (string, error) {
	token, err := r.Cookie("token")
	if err != nil {
		return "", err

	}
	return token.Value, nil
}
