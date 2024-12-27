package middleware

import (
	"backend/internal/utils"
	"context"
	"net/http"
)

// Add at package level
type contextKey string

const JwtTokenKey contextKey = "jwt_token"

func CheckJwtToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			utils.JSONError(rw, "Missing Authorization header", http.StatusUnauthorized)
			return
		}
		token = token[7:] // Remove "Bearer " prefix

		ctx := context.WithValue(r.Context(), JwtTokenKey, token)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
