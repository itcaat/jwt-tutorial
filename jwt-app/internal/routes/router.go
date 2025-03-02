package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/itcaat/jwt-tutorial/jwt-app/internal/handlers"
	"github.com/itcaat/jwt-tutorial/jwt-app/internal/middleware"
)

// RegisterRoutes добавляет маршруты в роутер
func RegisterRoutes(r *mux.Router) {
	// Обработка CORS для всех маршрутов
	r.Use(middleware.CORSMiddleware)

	// Публичные маршруты
	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")

	// Обработка preflight-запросов OPTIONS для всех путей
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			log.Println("OPTIONS request received")
			w.WriteHeader(http.StatusOK)
			return
		}
	}).Methods("OPTIONS")

	// Защищённые маршруты
	protected := r.PathPrefix("/protected").Subrouter()
	protected.Use(middleware.JWTMiddleware)
	protected.HandleFunc("/read", handlers.ReadHandler).Methods("GET")
	protected.HandleFunc("/write", handlers.WriteHandler).Methods("POST")
}
