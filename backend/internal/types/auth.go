package types

import "github.com/golang-jwt/jwt/v5"

type AuthUser struct {
	Email string `json:"email"`
	ID    string `json:"sub"`
	jwt.RegisteredClaims
}
