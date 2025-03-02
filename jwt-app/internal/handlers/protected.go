package handlers

import (
	"encoding/json"
	"net/http"
)

// ReadHandler – доступен всем авторизованным пользователям
func ReadHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("X-Username")
	role := r.Header.Get("X-Role")

	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Protected data accessed!",
		"username": username,
		"role":     role,
	})
}

// WriteHandler – доступен только для пользователей с ролью "write"
func WriteHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("X-Username")
	role := r.Header.Get("X-Role")

	if role != "write" {
		http.Error(w, "Forbidden: requires write access", http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Data successfully written!",
		"username": username,
	})
}
