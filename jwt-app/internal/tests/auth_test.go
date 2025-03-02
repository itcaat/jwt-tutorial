package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/itcaat/jwt-tutorial/jwt-app/internal/handlers"
	"github.com/itcaat/jwt-tutorial/jwt-app/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestLoginHandler(t *testing.T) {
	// Тестовые данные
	requestBody, _ := json.Marshal(models.Credentials{
		Username: "reader",
		Password: "password",
	})

	// Создаём HTTP-запрос
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Создаём тестовый HTTP-рекордер
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.LoginHandler)
	handler.ServeHTTP(rr, req)

	// Проверяем код ответа 200
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем, что в ответе есть токен
	var resp map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp, "token")
}
