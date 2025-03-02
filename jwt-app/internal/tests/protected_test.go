package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/itcaat/jwt-tutorial/jwt-app/internal/handlers"
	"github.com/stretchr/testify/assert"
)

func TestReadHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/protected/read", nil)
	req.Header.Set("X-Username", "reader")
	req.Header.Set("X-Role", "read")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.ReadHandler)
	handler.ServeHTTP(rr, req)

	// Ожидаем 200 OK
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestWriteHandler_ForbiddenForReader(t *testing.T) {
	req, _ := http.NewRequest("POST", "/protected/write", nil)
	req.Header.Set("X-Username", "reader")
	req.Header.Set("X-Role", "read") // У "reader" нет прав на запись

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.WriteHandler)
	handler.ServeHTTP(rr, req)

	// Ожидаем 403 Forbidden
	assert.Equal(t, http.StatusForbidden, rr.Code)
}
