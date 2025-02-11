package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/mahabub618/minipack/internal/models"
	"github.com/mahabub618/minipack/internal/services"
)

type PlatformHandler struct {
	platformService *services.PlatformService
}

func NewPlatformHandler(platformService *services.PlatformService) *PlatformHandler {
	return &PlatformHandler{platformService: platformService}
}

// CreatePlatform creates a new platform
func (h *PlatformHandler) CreatePlatform(w http.ResponseWriter, r *http.Request) {
	var platform models.Platform

	// Decode reqeust body into platform
	if err := json.NewDecoder(r.Body).Decode(&platform); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// validate reqeust body
	validate := validator.New()
	if err := validate.Struct(platform); err != nil {
		http.Error(w, "validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Create platform
	if err := h.platformService.CreatePlatform(r.Context(), &platform); err != nil {
		http.Error(w, "Failed to create platform: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the created platform
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(platform)
}

// GetPlatformByID retrieves a platform by its ID
func (h *PlatformHandler) GetPlatformByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid platform ID", http.StatusBadRequest)
		return
	}
	// Get platform by ID
	platform, err := h.platformService.GetPlatformByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to get platform: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Respond with the platform
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(platform)
}

// UpdatePlatform updates a platform
func (h *PlatformHandler) UpdatePlatform(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid platform ID", http.StatusBadRequest)
		return
	}
	var platform models.Platform
	// Decode request body into platform
	if err := json.NewDecoder(r.Body).Decode(&platform); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Validate request body
	validate := validator.New()
	if err := validate.Struct(platform); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Set the platform ID

	platform.ID = id

	// Update platform
	if err := h.platformService.UpdatePlatform(r.Context(), &platform); err != nil {
		http.Error(w, "Failed to update platform: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Respond with the updated platform
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(platform)
}

// DeletePlatform deletes a platform
func (h *PlatformHandler) DeletePlatform(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid platform ID", http.StatusBadRequest)
		return
	}
	// Delete platform
	if err := h.platformService.DeletePlatform(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete platform: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Respond with success
	w.WriteHeader(http.StatusNoContent)
}

// ListPlatforms retrieves a list of platforms
func (h *PlatformHandler) ListPlatforms(w http.ResponseWriter, r *http.Request) {
	// List platforms
	platforms, err := h.platformService.ListPlatforms(r.Context())
	if err != nil {
		http.Error(w, "Failed to list platforms: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Respond with the platforms
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(platforms)
}

// ActivatePlatform activates a platform
func (h *PlatformHandler) ActivatePlatform(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid platform ID", http.StatusBadRequest)
		return
	}
	// Activate platform
	if err := h.platformService.ActivatePlatform(r.Context(), id); err != nil {
		http.Error(w, "Failed to activate platform: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Respond with success
	w.WriteHeader(http.StatusNoContent)
}

// DeactivatePlatform deactivates a platform
func (h *PlatformHandler) DeactivatePlatform(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid platform ID", http.StatusBadRequest)
		return
	}
	// Deactivate platform
	if err := h.platformService.DeactivatePlatform(r.Context(), id); err != nil {
		http.Error(w, "Failed to deactivate platform: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Respond with success
	w.WriteHeader(http.StatusNoContent)
}

// PlatformExists checks if a platform exists
func (h *PlatformHandler) PlatformExists(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid platform ID", http.StatusBadRequest)
		return
	}
	// Check if platform exists
	exists, err := h.platformService.GetPlatformByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to check platform: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Respond with exists
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(exists)
}
