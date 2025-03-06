package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mahabub618/minipack/internal/services"
)

type PackageHandler struct {
	packageService *services.PackageService
}

func NewPackageHandler(packageService *services.PackageService) *PackageHandler {
	return &PackageHandler{packageService: packageService}
}

// GetPackageByID retrieves a package by its ID
func (h *PackageHandler) GetPackageByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid package ID", http.StatusBadRequest)
		return
	}
	// Get package by ID and check if it exists
	pkg, err := h.packageService.GetPackageByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to get package: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if pkg == nil {
		http.Error(w, "Package not found", http.StatusNotFound)
		return
	}

	// Respond with the package
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pkg)
}

// ListPackagesByPlatform retrieves all packages for a specific platform.
func (h *PackageHandler) ListPackagesByPlatform(w http.ResponseWriter, r *http.Request) {
	platformIDStr := chi.URLParam(r, "id")
	platformID, err := strconv.Atoi(platformIDStr)
	if err != nil {
		http.Error(w, "Invalid platform ID", http.StatusBadRequest)
		return
	}

	// Retrieve packages for the platform
	packages, err := h.packageService.ListPackagesByPlatform(r.Context(), platformID)
	if err != nil {
		http.Error(w, "Failed to retrieve packages: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if packages == nil {
		http.Error(w, "No packages found", http.StatusNotFound)
		return
	}

	// Respond with the packages
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(packages)
}
