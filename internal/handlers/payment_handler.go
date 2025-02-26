package handlers

import (
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/mahabub618/minipack/config"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/mahabub618/minipack/internal/models"
	"github.com/mahabub618/minipack/internal/services"
)

type PaymentHandler struct {
	paymentService      *services.PaymentService
	subscriptionService *services.SubscriptionService
}

func NewPaymentHandler(paymentService *services.PaymentService, subscriptionService *services.SubscriptionService) *PaymentHandler {
	return &PaymentHandler{
		paymentService:      paymentService,
		subscriptionService: subscriptionService,
	}
}

// CreatePayment handles the creation of a new payment.
func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var payment models.Payment

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request body
	validate := validator.New()
	if err := validate.Struct(payment); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get the subscription from SubscriptionService
	subscription, err := h.subscriptionService.GetSubscriptionByID(r.Context(), payment.SubscriptionID)
	if err != nil {
		http.Error(w, "Failed to retrieve subscription: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if subscription == nil {
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}

	// Authenticate and authorize user
	unauthorized := h.authenticateAndAuthorize(r, subscription.UserID)
	if unauthorized != nil {
		http.Error(w, "Forbidden", http.StatusUnauthorized)
		return
	}

	// Create the payment
	if err := h.paymentService.CreatePayment(r.Context(), &payment); err != nil {
		http.Error(w, "Failed to create payment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the created payment
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payment)
}

func (h *PaymentHandler) authenticateAndAuthorize(r *http.Request, userID int) error {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return errors.New("unauthorized")
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	cfg, err := config.LoadConfig()
	if err != nil {
		return errors.New("internal Server Error: Unable to load config")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return errors.New("unauthorized")
	}

	if claims.UserID != userID {
		return errors.New("forbidden")
	}

	return nil
}

// ConfirmPayment confirms a payment.
func (h *PaymentHandler) ConfirmPayment(w http.ResponseWriter, r *http.Request) {
	paymentIDStr := chi.URLParam(r, "id")
	paymentID, err := strconv.Atoi(paymentIDStr)
	if err != nil {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}

	// Get the payment
	payment, err := h.paymentService.GetPaymentByID(r.Context(), paymentID)
	if err != nil {
		http.Error(w, "Failed to retrieve payment: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if payment == nil {
		http.Error(w, "Payment not found", http.StatusNotFound)
		return
	}

	// Get the subscription
	subscription, err := h.subscriptionService.GetSubscriptionByID(r.Context(), payment.SubscriptionID)
	if err != nil {
		http.Error(w, "Failed to retrieve subscription: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if subscription == nil {
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}

	// Authenticate and authorize user
	unauthorized := h.authenticateAndAuthorize(r, subscription.UserID)
	if unauthorized != nil {
		http.Error(w, "Forbidden", http.StatusUnauthorized)
		return
	}

	// Confirm the payment
	if err := h.paymentService.ConfirmPayment(r.Context(), paymentID); err != nil {
		http.Error(w, "Failed to confirm payment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
}

// GetPaymentByID retrieves a payment by its ID.
func (h *PaymentHandler) GetPaymentByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}

	// Retrieve the payment
	payment, err := h.paymentService.GetPaymentByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to retrieve payment: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if payment == nil {
		http.Error(w, "Payment not found", http.StatusNotFound)
		return
	}

	// Respond with the payment
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payment)
}

// ListPaymentsBySubscription retrieves all payments for a specific subscription.
func (h *PaymentHandler) ListPaymentsBySubscription(w http.ResponseWriter, r *http.Request) {
	subscriptionIDStr := chi.URLParam(r, "subscription_id")
	subscriptionID, err := strconv.Atoi(subscriptionIDStr)
	if err != nil {
		http.Error(w, "Invalid subscription ID", http.StatusBadRequest)
		return
	}

	// Retrieve payments for the subscription
	payments, err := h.paymentService.ListPaymentsBySubscription(r.Context(), subscriptionID)
	if err != nil {
		http.Error(w, "Failed to retrieve payments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the payments
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payments)
}
