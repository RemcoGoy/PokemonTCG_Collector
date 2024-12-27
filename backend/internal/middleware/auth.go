package middleware

import (
	"backend/internal/types"
	"backend/internal/utils"
	"context"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// Add at package level
type contextKey string

const JwtTokenKey contextKey = "jwt_token"
const UserEmail contextKey = "user_email"
const UserID contextKey = "user_id"

func CheckJwtToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			utils.JSONError(rw, "Missing Authorization header", http.StatusUnauthorized)
			return
		}
		tokenString = tokenString[7:] // Remove "Bearer " prefix

		token, err := jwt.ParseWithClaims(tokenString, &types.AuthUser{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			utils.JSONError(rw, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*types.AuthUser)
		if !ok {
			utils.JSONError(rw, "unknown claims type, cannot proceed", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), JwtTokenKey, tokenString)
		ctx = context.WithValue(ctx, UserEmail, claims.Email)
		ctx = context.WithValue(ctx, UserID, claims.ID)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
