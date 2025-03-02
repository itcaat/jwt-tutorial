package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/itcaat/jwt-tutorial/jwt-app/internal/auth"
	"github.com/itcaat/jwt-tutorial/jwt-app/internal/config"
	"github.com/itcaat/jwt-tutorial/jwt-app/internal/models"
)

// LoginHandler – аутентификация и выдача JWT
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, exists := auth.Users[creds.Username] // <-- теперь используем auth.Users
	if !exists || user.Password != creds.Password {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Создаём JWT-токен
	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &models.Claims{
		Username: creds.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
