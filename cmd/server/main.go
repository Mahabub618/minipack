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
	packageRepo := repositories.NewPackageRepository(db)
	subscriptionRepo := repositories.NewSubscriptionRepository(db)
	paymentRepo := repositories.NewPaymentRepository(db)

	userService := services.NewUserService(userRepo)
	platformService := services.NewPlatformService(platformRepo)
	packageService := services.NewPackageService(packageRepo, platformRepo)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo, packageRepo)
	paymentService := services.NewPaymentService(paymentRepo)

	userHandler := handlers.NewUserHandler(userService)
	platformHandler := handlers.NewPlatformHandler(platformService)
	packageHandler := handlers.NewPackageHandler(packageService)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService, userService)
	paymentHandler := handlers.NewPaymentHandler(paymentService, subscriptionService)

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Public Routes
	router.Get("/ping", handlers.PingHandler)
	router.Post("/auth/signup", userHandler.RegisterHandler)
	router.Post("/auth/login", userHandler.LoginHandler)

	// User Routes
	router.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware("user"))
		r.Get("/platforms", platformHandler.ListPlatforms)
		r.Get("/platforms/{id}", platformHandler.GetPlatformByID)

		r.Get("/packages/{id}", packageHandler.GetPackageByID)
		r.Get("/packages/list/{id}", packageHandler.ListPackagesByPlatform)

		r.Post("/subscriptions", subscriptionHandler.CreateSubscription)

		r.Post("/payments", paymentHandler.CreatePayment)
		r.Get("/payments/{id}", paymentHandler.GetPaymentByID)
		r.Get("/payments/subscription/{subscription_id}", paymentHandler.ListPaymentsBySubscription)
	})

	// Admin Routes
	router.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware("admin"))
		r.Post("/platforms", platformHandler.CreatePlatform)
		r.Put("/platforms/{id}", platformHandler.UpdatePlatform)
		r.Delete("/platforms/{id}", platformHandler.DeletePlatform)
		r.Post("/platforms/{id}/activate", platformHandler.ActivatePlatform)
		r.Post("/platforms/{id}/deactivate", platformHandler.DeactivatePlatform)

		r.Post("/packages", packageHandler.CreatePackage)
		r.Put("/packages/{id}", packageHandler.UpdatePackage)
		r.Delete("/packages/{id}", packageHandler.DeletePackage)

		r.Post("/payments/{id}/confirm", paymentHandler.ConfirmPayment)
	})

	log.Println("Starting server on: 8585..")
	if err := http.ListenAndServe(":8585", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
