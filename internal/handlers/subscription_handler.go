package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/mahabub618/minipack/internal/models"
	"github.com/mahabub618/minipack/internal/services"
)

type SubscriptionHandler struct {
	subscriptionService *services.SubscriptionService
	userService         *services.UserService
	packageService      *services.PackageService
}

func NewSubscriptionHandler(subscriptionService *services.SubscriptionService, userService *services.UserService, packageService *services.PackageService) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
		userService:         userService,
		packageService:      packageService,
	}
}

// CreateSubscription handles the creation of a new subscription.
func (h *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var sub models.Subscription

	// Decode request body with a limit
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1048576))
	if err := decoder.Decode(&sub); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Check if user exists
	if exists, err := h.userService.UserExists(r.Context(), sub.UserID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to check user")
		return
	} else if !exists {
		respondWithError(w, http.StatusBadRequest, "User does not exist")
		return
	}

	// Fetch package details
	pkg, err := h.packageService.GetPackageByID(r.Context(), sub.PackageID)
	if err != nil || pkg == nil {
		respondWithError(w, http.StatusBadRequest, "Package does not exist")
		return
	}

	// Set subscription details
	now := time.Now()
	sub.CreatedAt, sub.UpdatedAt, sub.StartDate = now, now, now
	sub.EndDate = now.Add(time.Duration(pkg.Duration) * 24 * time.Hour)
	sub.Price, sub.Currency = pkg.Price-pkg.DiscountAmount, pkg.Currency
	sub.Status = "active"

	// Validate request
	if err := validate.Struct(sub); err != nil {
		respondWithError(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	// Create subscription
	if err := h.subscriptionService.CreateSubscription(r.Context(), &sub); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create subscription")
		return
	}

	// Respond with only the "status" field
	response := struct {
		Status string `json:"status"`
	}{Status: sub.Status}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetSubscriptionByID retrieves a subscription by its ID.
func (h *SubscriptionHandler) GetSubscriptionByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid subscription ID", http.StatusBadRequest)
		return
	}

	// Retrieve the subscription
	sub, err := h.subscriptionService.GetSubscriptionByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to retrieve subscription: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if sub == nil {
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}

	// Respond with the subscription
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sub)
}

// UpdateSubscription updates a subscription's details.
func (h *SubscriptionHandler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid subscription ID", http.StatusBadRequest)
		return
	}

	var sub map[string]interface{}

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Get existing subscription
	existingSub, err := h.subscriptionService.GetSubscriptionByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to get subscription: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if existingSub == nil {
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}

	// Merge existing subscription with request body
	for key, value := range sub {
		switch key {
		case "user_id":
			existingSub.UserID = value.(int)
		case "package_id":
			existingSub.PackageID = value.(int)
		case "start_date":
			existingSub.StartDate = value.(time.Time)
		case "end_date":
			existingSub.EndDate = value.(time.Time)
		case "status":
			existingSub.Status = value.(string)
		case "auto_renew":
			existingSub.AutoRenew = value.(bool)
		case "trial_end_date":
			existingSub.TrialEndDate = value.(time.Time)
		case "next_billing_date":
			existingSub.NextBillingDate = value.(time.Time)
		}
	}

	// validate the request body
	validate := validator.New()
	if err := validate.Struct(sub); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Update the subscription
	if err := h.subscriptionService.UpdateSubscription(r.Context(), existingSub); err != nil {
		http.Error(w, "Failed to update subscription: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the updated subscription
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sub)
}

// DeleteSubscription deletes a subscription by its ID.
func (h *SubscriptionHandler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid subscription ID", http.StatusBadRequest)
		return
	}

	// Delete the subscription
	if err := h.subscriptionService.DeleteSubscription(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete subscription: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusNoContent)
}

// ListSubscriptionsByUser retrieves all subscriptions for a specific user.
func (h *SubscriptionHandler) ListSubscriptionsByUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Retrieve subscriptions for the user
	subscriptions, err := h.subscriptionService.ListSubscriptionsByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to retrieve subscriptions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the subscriptions
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(subscriptions)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
