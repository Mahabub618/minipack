package services

import (
	"context"
	"errors"
	"github.com/mahabub618/minipack/internal/models"
	"github.com/mahabub618/minipack/internal/repositories"
)

var ErrPackageExists = errors.New("package name already exists")
var ErrPackageNotFound = errors.New("package not found")

type PackageService struct {
	packageRepo  *repositories.PackageRepository
	platformRepo *repositories.PlatformRepository
}

func NewPackageService(packageRepo *repositories.PackageRepository, platformRepo *repositories.PlatformRepository) *PackageService {
	return &PackageService{
		packageRepo:  packageRepo,
		platformRepo: platformRepo,
	}
}

// GetPackageByID retrieves a package by its ID
func (s *PackageService) GetPackageByID(ctx context.Context, id int) (*models.Package, error) {
	// Check if the package exists in the database
	return s.packageRepo.FindPackageByID(ctx, id)
}

// ListPackagesByPlatform retrieves all packages for a specific platform.
func (s *PackageService) ListPackagesByPlatform(ctx context.Context, platformID int) ([]models.Package, error) {
	return s.packageRepo.ListPackagesByPlatform(ctx, platformID)
}
