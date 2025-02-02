package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/mahabub618/minipack/config"
	"github.com/mahabub618/minipack/internal/database"
	"github.com/mahabub618/minipack/internal/handlers"
	"log"
	"net/http"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	if err := database.InitDB(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/ping", handlers.PingHandler)
	router.Post("/auth/signup", handlers.SignupHandler)

	log.Println("Starting server on: 8585..")
	if err := http.ListenAndServe(":8585", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
