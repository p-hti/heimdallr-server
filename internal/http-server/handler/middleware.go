package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/p-hti/heimdallr-server/internal/domain/model"
)

const userContextKey model.ContextKey = "UserID"

func (h *HTTPServer) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getTokenFromRequest(r)
		if err != nil {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		userID, err := h.Service.ParseToken(r.Context(), token)
		if err != nil {
			http.Error(w, "failed to parse token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getTokenFromRequest(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("missing authorization header")
	}

	tokenString := strings.TrimPrefix(header, "Bearer ")
	if tokenString == header {
		return "", errors.New("invalid authorization format")
	}
	return tokenString, nil
}
