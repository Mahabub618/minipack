package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/mahabub618/minipack/internal/handlers"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/ping", handlers.PingHandler)

	log.Println("Starting server on: 8585..")
	if err := http.ListenAndServe(":8585", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
