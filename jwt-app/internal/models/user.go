package models

import "github.com/golang-jwt/jwt/v5"

// Claims – структура данных, хранящаяся в JWT-токене
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}
