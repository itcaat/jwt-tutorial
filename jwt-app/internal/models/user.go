package models

import "github.com/golang-jwt/jwt/v5"

// Claims – структура access-токена
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// RefreshClaims – структура refresh-токена
type RefreshClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}
