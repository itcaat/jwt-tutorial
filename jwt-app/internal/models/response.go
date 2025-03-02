package models

// AuthResponse – ответ при успешной аутентификации
type AuthResponse struct {
	Token string `json:"token"`
}

// ErrorResponse – стандартный формат ошибки
type ErrorResponse struct {
	Error string `json:"error"`
}
