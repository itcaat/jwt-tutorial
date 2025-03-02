package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/itcaat/jwt-tutorial/jwt-app/internal/middleware"
	"github.com/stretchr/testify/assert"
)

func TestJWTMiddleware_Unauthorized(t *testing.T) {
	req, _ := http.NewRequest("GET", "/protected/read", nil)
	rr := httptest.NewRecorder()

	handler := middleware.JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	// Ожидаем 401 Unauthorized, так как токена нет
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
