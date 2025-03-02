package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/itcaat/jwt-tutorial/jwt-app/internal/config"
	"github.com/itcaat/jwt-tutorial/jwt-app/internal/models"
)

// JWTMiddleware проверяет JWT-токен в заголовке Authorization
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &models.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JWTSecret), nil // <-- Используем секрет из env
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Передаём username и роль в заголовки запроса
		r.Header.Set("X-Username", claims.Username)
		r.Header.Set("X-Role", claims.Role)

		next.ServeHTTP(w, r)
	})
}
