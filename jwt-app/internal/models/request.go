package models

// Credentials – структура для входа (логин и пароль)
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
