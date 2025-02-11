package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/mahabub618/minipack/config"
	"github.com/mahabub618/minipack/internal/database"
	"github.com/mahabub618/minipack/internal/handlers"
	"github.com/mahabub618/minipack/internal/repositories"
	"github.com/mahabub618/minipack/internal/services"
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

	db := database.GetDB()
	userRepo := repositories.NewUserRepository(db)
	platformRepo := repositories.NewPlatformRepository(db)

	userService := services.NewUserService(userRepo)
	platformService := services.NewPlatformService(platformRepo)

	userHandler := handlers.NewUserHandler(userService)
	platformHandler := handlers.NewPlatformHandler(platformService)

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/ping", handlers.PingHandler)
	router.Post("/auth/signup", userHandler.RegisterHandler)
	router.Post("/auth/login", userHandler.LoginHandler)

	router.Post("/platforms", platformHandler.CreatePlatform)
	router.Get("/platforms/{id}", platformHandler.GetPlatformByID)
	router.Put("/platforms/{id}", platformHandler.UpdatePlatform)
	router.Delete("/platforms/{id}", platformHandler.DeletePlatform)
	router.Get("/platforms", platformHandler.ListPlatforms)
	router.Post("/platforms/{id}/activate", platformHandler.ActivatePlatform)
	router.Post("/platforms/{id}/deactivate", platformHandler.DeactivatePlatform)

	log.Println("Starting server on: 8585..")
	if err := http.ListenAndServe(":8585", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
