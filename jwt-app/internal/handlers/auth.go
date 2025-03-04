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

// LoginHandler – выдаёт access и refresh токены
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, exists := auth.Users[creds.Username]
	if !exists || user.Password != creds.Password {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Генерируем Access Token (JWT)
	accessToken, err := generateAccessToken(creds.Username, user.Role)
	if err != nil {
		http.Error(w, "Could not generate access token", http.StatusInternalServerError)
		return
	}

	// Генерируем Refresh Token (JWT)
	refreshToken, err := generateRefreshToken(creds.Username)
	if err != nil {
		http.Error(w, "Could not generate refresh token", http.StatusInternalServerError)
		return
	}

	// Устанавливаем refresh-token в HttpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour), // 7 дней
		Path:     "/refresh",
	})

	// Отправляем access_token в JSON-ответе
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
	})
}

// RefreshHandler – обновляет access-токен по refresh-токену
func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Refresh token required", http.StatusUnauthorized)
		return
	}

	refreshToken := cookie.Value

	// Проверяем refresh-токен
	claims, err := validateRefreshToken(refreshToken)
	if err != nil {
		http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	// Генерируем новый access-токен
	accessToken, err := generateAccessToken(claims.Username, claims.Role)
	if err != nil {
		http.Error(w, "Could not generate new access token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
	})
}

// generateAccessToken создаёт короткоживущий JWT (10 мин)
func generateAccessToken(username, role string) (string, error) {
	expirationTime := time.Now().Add(10)

	claims := &models.Claims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWTSecret))
}

// generateRefreshToken создаёт refresh-токен (7 дней)
func generateRefreshToken(username string) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour) // 7 дней

	claims := &models.RefreshClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWTSecret))
}

// validateRefreshToken проверяет refresh-токен
func validateRefreshToken(refreshToken string) (*models.RefreshClaims, error) {
	claims := &models.RefreshClaims{}

	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
