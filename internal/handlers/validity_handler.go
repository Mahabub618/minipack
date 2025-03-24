package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/mahabub618/minipack/internal/models"
	"github.com/mahabub618/minipack/internal/services"
)

type ValidityHandler struct {
	validityServices *services.ValidityService
	platformServices *services.PlatformService
}

func NewValidityHandler(validityServices *services.ValidityService, platformServices *services.PlatformService) *ValidityHandler {
	return &ValidityHandler{
		validityServices: validityServices,
		platformServices: platformServices,
	}
}

func (h *ValidityHandler) CreateValidity(w http.ResponseWriter, r *http.Request) {
	var validity models.Validity

	if err := json.NewDecoder(r.Body).Decode(&validity); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(validity); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	platform, err := h.platformServices.GetPlatformByID(r.Context(), validity.PlatformID)
	if err != nil {
		http.Error(w, "Failed to get platform", http.StatusInternalServerError)
		return
	}
	if platform == nil {
		http.Error(w, "Platform not found", http.StatusBadRequest)
		return
	}

	if err := h.validityServices.CreateValidity(r.Context(), &validity); err != nil {
		http.Error(w, "Failed to create validity: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(validity)
}

func (h *ValidityHandler) UpdateValidity(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Validity ID", http.StatusBadRequest)
		return
	}

	var validity map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&validity); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	existingValidity, err := h.validityServices.GetValidityByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to get validity: "+err.Error(), http.StatusBadRequest)
		return
	}
	if existingValidity == nil {
		http.Error(w, "Validity not found", http.StatusNotFound)
		return
	}

	for key, value := range validity {
		switch key {
		case "platform_id":
			existingValidity.PlatformID = value.(int)
		case "duration":
			existingValidity.Duration = value.(int)
		case "price":
			existingValidity.Price = value.(float64)
		case "label":
			existingValidity.Label = value.(string)
		}
	}
	existingValidity.UpdatedAt = time.Now()

	validate := validator.New()
	if err := validate.Struct(existingValidity); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.validityServices.UpdateValidity(r.Context(), existingValidity); err != nil {
		http.Error(w, "Failed to get validity: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingValidity)
}

func (h *ValidityHandler) GetValidityById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Validity ID", http.StatusBadRequest)
		return
	}

	validity, err := h.validityServices.GetValidityByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Validity id not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get validity: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if validity == nil {
		http.Error(w, "Validity not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(validity)
}

func (h *ValidityHandler) DeleteValidity(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Validity ID", http.StatusBadRequest)
		return
	}

	if err := h.validityServices.DeleteValidityByID(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete validity", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ValidityHandler) GetAllValiditiesByPlatformId(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Platform ID", http.StatusBadRequest)
		return
	}

	platform, err := h.platformServices.GetPlatformByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to get platform", http.StatusInternalServerError)
		return
	}
	if platform == nil {
		http.Error(w, "Platform not found", http.StatusBadRequest)
		return
	}

	validityList, err := h.validityServices.GetAllValiditiesByPlatformID(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to list validity: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(validityList)
}
