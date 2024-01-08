package transport

import (
	"context"
	"net/http"

	"github.com/AlexCorn999/transaction-system/internal/domain"
	"github.com/AlexCorn999/transaction-system/internal/logger"
)

// authMiddleware performs the middleware authorization function.
// Gets the token from the request and passes to the context the userID that makes the request.
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

// getTokenFromRequest gets the token from the cookie.
func getTokenFromRequest(r *http.Request) (string, error) {
	token, err := r.Cookie("token")
	if err != nil {
		return "", err

	}
	return token.Value, nil
}
