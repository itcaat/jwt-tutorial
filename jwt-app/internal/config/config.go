package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var JWTSecret string

func LoadConfig() {
	// Загружаем переменные из .env, если файл существует
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	JWTSecret = os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}
}
