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

type PackageHandler struct {
	packageService *services.PackageService
}

func NewPackageHandler(packageService *services.PackageService) *PackageHandler {
	return &PackageHandler{packageService: packageService}
}

// CreatePackage creates a new package
func (h *PackageHandler) CreatePackage(w http.ResponseWriter, r *http.Request) {
	var pkg models.Package

	// Decode reqeust body into package
	if err := json.NewDecoder(r.Body).Decode(&pkg); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// validate reqeust body
	validate := validator.New()
	if err := validate.Struct(pkg); err != nil {
		http.Error(w, "validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Create package
	if err := h.packageService.CreatePackage(r.Context(), &pkg); err != nil {
		http.Error(w, "Failed to create package: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the created package
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pkg)
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

// UpdatePackage updates a package
func (h *PackageHandler) UpdatePackage(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid package ID", http.StatusBadRequest)
		return
	}

	var pkg models.Package
	// Decode reqeust body into package
	if err := json.NewDecoder(r.Body).Decode(&pkg); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	pkg.ID = id
	// validate reqeust body
	validate := validator.New()
	if err := validate.Struct(pkg); err != nil {
		http.Error(w, "validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Update package
	if err := h.packageService.UpdatePackage(r.Context(), &pkg); err != nil {
		http.Error(w, "Failed to update package: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the updated package
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pkg)
}

// DeletePackage deletes a package
func (h *PackageHandler) DeletePackage(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid package ID", http.StatusBadRequest)
		return
	}
	// Delete package
	if err := h.packageService.DeletePackage(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete package: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
}

// ListPackagesByPlatform retrieves all packages for a specific platform.
func (h *PackageHandler) ListPackagesByPlatform(w http.ResponseWriter, r *http.Request) {
	platformIDStr := chi.URLParam(r, "platform_id")
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
