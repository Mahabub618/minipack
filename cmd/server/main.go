package main

import (
	"log"
	"net/http"

	"github.com/mahabub618/minipack/internal/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/mahabub618/minipack/config"
	"github.com/mahabub618/minipack/internal/database"
	"github.com/mahabub618/minipack/internal/handlers"
	"github.com/mahabub618/minipack/internal/repositories"
	"github.com/mahabub618/minipack/internal/services"
)

func main() {
	//err := godotenv.Load(".env")
	//if err != nil {
	//	log.Fatalf("Error loading .env file: %v", err)
	//}

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
	validityRepo := repositories.NewValidityRepository(db)
	subscriptionRepo := repositories.NewSubscriptionRepository(db)
	paymentRepo := repositories.NewPaymentRepository(db)

	userService := services.NewUserService(userRepo)
	platformService := services.NewPlatformService(platformRepo)
	packageService := services.NewPackageService(packageRepo, platformRepo)
	validityServices := services.NewValidityService(validityRepo)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo, validityRepo)
	paymentService := services.NewPaymentService(paymentRepo)

	userHandler := handlers.NewUserHandler(userService)
	platformHandler := handlers.NewPlatformHandler(platformService)
	packageHandler := handlers.NewPackageHandler(packageService, platformService, validityServices)
	validityHander := handlers.NewValidityHandler(validityServices, platformService)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService, userService)
	paymentHandler := handlers.NewPaymentHandler(paymentService, subscriptionService)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"}, // Allow Angular app
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not to send preflight requests
	}))

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Public Routes
	router.Get("/ping", handlers.PingHandler)
	router.Post("/auth/signup", userHandler.RegisterHandler)
	router.Post("/auth/login", userHandler.LoginHandler)

	router.Get("/platforms", platformHandler.ListPlatforms)
	router.Get("/platforms/{id}", platformHandler.GetPlatformByID)

	router.Get("/validity/{id}", validityHander.GetValidityById)
	router.Get("/validity/list/{id}", validityHander.GetAllValiditiesByPlatformId)

	router.Get("/packages/{id}", packageHandler.GetPackageValidityByID)
	router.Get("/packages/list/{id}", packageHandler.ListPackagesByPlatform)

	router.Post("/subscriptions", subscriptionHandler.CreateSubscription)

	router.Post("/payments", paymentHandler.CreatePayment)
	router.Get("/payments/{id}", paymentHandler.GetPaymentByID)
	router.Get("/payments/subscription/{subscription_id}", paymentHandler.ListPaymentsBySubscription)

	// User Routes
	//router.Group(func(r chi.Router) {
	//	r.Use(middlewares.AuthMiddleware("user", "admin"))
	//	r.Get("/platforms", platformHandler.ListPlatforms)
	//	r.Get("/platforms/{id}", platformHandler.GetPlatformByID)
	//
	//	r.Get("/validity/{id}", validityHander.GetValidityById)
	//	r.Get("/validity/list/{id}", validityHander.GetAllValiditiesByPlatformId)
	//
	//	r.Get("/packages/{id}", packageHandler.GetPackageValidityByID)
	//	r.Get("/packages/list/{id}", packageHandler.ListPackagesByPlatform)
	//
	//	r.Post("/subscriptions", subscriptionHandler.CreateSubscription)
	//
	//	r.Post("/payments", paymentHandler.CreatePayment)
	//	r.Get("/payments/{id}", paymentHandler.GetPaymentByID)
	//	r.Get("/payments/subscription/{subscription_id}", paymentHandler.ListPaymentsBySubscription)
	//})

	// Admin Routes
	router.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware("admin"))
		r.Post("/platforms", platformHandler.CreatePlatform)
		r.Put("/platforms/{id}", platformHandler.UpdatePlatform)
		r.Delete("/platforms/{id}", platformHandler.DeletePlatform)
		r.Post("/platforms/{id}/activate", platformHandler.ActivatePlatform)
		r.Post("/platforms/{id}/deactivate", platformHandler.DeactivatePlatform)

		r.Post("/validity", validityHander.CreateValidity)
		r.Put("/validity/{id}", validityHander.UpdateValidity)
		r.Delete("/validity/{id}", validityHander.DeleteValidity)

		r.Post("/payments/{id}/confirm", paymentHandler.ConfirmPayment)
	})

	log.Println("Starting server on: 8585..")
	if err := http.ListenAndServe(":8585", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
