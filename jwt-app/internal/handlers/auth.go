package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		Domain:   "localhost.devopsbrain.ru",
	})

	// Отправляем token в JSON-ответе
	json.NewEncoder(w).Encode(map[string]string{
		"token": accessToken,
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
	username, err := validateRefreshToken(refreshToken)
	if err != nil {
		http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	// ✅ 1. Генерируем новый access-токен
	accessToken, err := generateAccessToken(username, auth.Users[username].Role)
	if err != nil {
		http.Error(w, "Could not generate new access token", http.StatusInternalServerError)
		return
	}

	// ✅ 2. Генерируем новый refresh-токен
	newRefreshToken, err := generateRefreshToken(username)
	if err != nil {
		http.Error(w, "Could not generate new refresh token", http.StatusInternalServerError)
		return
	}

	// ✅ 3. Сохраняем новый refresh-токен (опционально в БД/Redis)
	// storeRefreshTokenInDB(claims.Username, newRefreshToken)

	// ✅ 4. Удаляем старый refresh-токен (если храним в БД/Redis)
	// revokeRefreshToken(refreshToken)

	// ✅ 5. Отправляем новый refresh-токен в HttpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour), // 7 дней
		Path:     "/",
		Domain:   "localhost.devopsbrain.ru",
	})

	// ✅ 6. Отправляем новый access-токен в JSON-ответе
	json.NewEncoder(w).Encode(map[string]string{
		"token": accessToken,
	})
}

// generateAccessToken создаёт короткоживущий JWT (10 мин)
func generateAccessToken(username, role string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Minute) // 1 минут

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
	claims := jwt.MapClaims{
		"jti": uuid.NewString(),
		"sub": username,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWTSecret))
}

func validateRefreshToken(refreshToken string) (string, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		log.Println(err)
		return "", err
	}

	// Приводим claims к MapClaims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("invalid token claims")
		return "", errors.New("invalid token claims")
	}

	// Проверяем `exp` (если `exp` истёк, токен недействителен)
	exp, ok := claims["exp"].(float64) // JWT `exp` обычно приходит как float64 (Unix timestamp)
	if !ok {
		log.Println("invalid exp field in token")
		return "", errors.New("invalid exp field in token")
	}

	if time.Now().Unix() > int64(exp) {
		log.Println("refresh token expired")
		return "", errors.New("refresh token expired")
	}

	// Получаем `sub` (Subject = Username)
	username, ok := claims["sub"].(string)
	if !ok {
		log.Println("invalid sub field in token")
		return "", errors.New("invalid sub field in token")
	}

	_, exists := auth.Users[username]
	if !exists {
		log.Println("user not exist anymore")
		return "", errors.New("user not exist anymore")
	}

	return username, nil
}

// LogoutHandler – удаляет refresh-токен при выходе
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Удаляем refresh-токен из cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Expires:  time.Now().Add(-1 * time.Hour), // 7 дней
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		Domain:   "localhost.devopsbrain.ru",
	})

	// Можно добавить refresh_token в чёрный список (например, Redis)
	// blacklistedTokens[refreshToken] = true

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout successful",
	})
}
