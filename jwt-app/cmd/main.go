package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/itcaat/jwt-tutorial/jwt-app/internal/routes"
)

func main() {
	r := mux.NewRouter()
	// Регистрируем маршруты
	routes.RegisterRoutes(r)

	port := ":8080"
	fmt.Println("Server running on http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, r))
}
