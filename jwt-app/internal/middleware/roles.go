package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

// RoleMiddleware проверяет, что у пользователя есть требуемая роль
func RoleMiddleware(requiredRole string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := r.Header.Get("X-Role")

			if requiredRole == "write" && role != "write" {
				http.Error(w, "Forbidden: requires write access", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
