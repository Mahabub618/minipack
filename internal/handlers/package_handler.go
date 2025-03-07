package handlers

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/mahabub618/minipack/internal/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mahabub618/minipack/internal/services"
)

type PackageHandler struct {
	packageService   *services.PackageService
	plaformServices  *services.PlatformService
	validityServices *services.ValidityService
}

func NewPackageHandler(packageService *services.PackageService, plaformServices *services.PlatformService, validityServices *services.ValidityService) *PackageHandler {
	return &PackageHandler{
		packageService:   packageService,
		plaformServices:  plaformServices,
		validityServices: validityServices,
	}
}

// GetPackageByID retrieves a package by its ID
func (h *PackageHandler) GetPackageValidityByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid validity ID", http.StatusBadRequest)
		return
	}

	var pkg models.Package
	if err := json.NewDecoder(r.Body).Decode(&pkg); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(pkg); err != nil {
		http.Error(w, "Validation falied", http.StatusBadRequest)
		return
	}

	pkgValidity, err := h.validityServices.GetValidityByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to get package validity: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the package
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pkgValidity)
}

// ListPackagesByPlatform retrieves all packages for a specific platform.
func (h *PackageHandler) ListPackagesByPlatform(w http.ResponseWriter, r *http.Request) {
	platformIDStr := chi.URLParam(r, "id")
	platformID, err := strconv.Atoi(platformIDStr)
	if err != nil {
		http.Error(w, "Invalid platform ID", http.StatusBadRequest)
		return
	}

	var pkg models.Package
	if err := json.NewDecoder(r.Body).Decode(&pkg); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(pkg); err != nil {
		http.Error(w, "Validation falied", http.StatusBadRequest)
		return
	}

	plaform, err := h.plaformServices.GetPlatformByID(r.Context(), platformID)
	if err != nil {
		http.Error(w, "Failed to get platform", http.StatusInternalServerError)
		return
	}
	if plaform == nil {
		http.Error(w, "Platform not found", http.StatusNotFound)
		return
	}

	validityLists, err := h.validityServices.GetAllValiditiesByPlatformID(r.Context(), platformID)
	if err != nil {
		http.Error(w, "Failed to get validity list using platform id: "+err.Error(), http.StatusInternalServerError)
		return
	}

	platformWithValidityList := struct {
		*models.Platform
		ValidityLists []*models.Validity `json:"validity"`
	}{
		Platform:      plaform,
		ValidityLists: validityLists,
	}

	// Respond with the packages
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(platformWithValidityList)
}
