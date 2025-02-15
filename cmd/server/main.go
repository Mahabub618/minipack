package main

import (
	"log"
	"net/http"

	"github.com/mahabub618/minipack/internal/middlewares"

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

	router.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware("admin"))
		r.Post("/platforms", platformHandler.CreatePlatform)
		r.Get("/platforms/{id}", platformHandler.GetPlatformByID)
		r.Put("/platforms/{id}", platformHandler.UpdatePlatform)
		r.Delete("/platforms/{id}", platformHandler.DeletePlatform)
		r.Get("/platforms", platformHandler.ListPlatforms)
		r.Post("/platforms/{id}/activate", platformHandler.ActivatePlatform)
		r.Post("/platforms/{id}/deactivate", platformHandler.DeactivatePlatform)
	})

	log.Println("Starting server on: 8585..")
	if err := http.ListenAndServe(":8585", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
