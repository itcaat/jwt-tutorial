package handlers

import (
	"encoding/json"
	"net/http"
)

// HomeHandler отвечает на публичные запросы
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "Welcome to the public home page!"})
}
