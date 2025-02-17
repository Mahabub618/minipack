package services

import (
	"context"
	"errors"
	"github.com/mahabub618/minipack/internal/models"
	"github.com/mahabub618/minipack/internal/repositories"
	"time"
)

var ErrPackageExists = errors.New("package name already exists")
var ErrPackageNotFound = errors.New("package not found")

type PackageService struct {
	packageRepo *repositories.PackageRepository
	platformRepo *repositories.PlatformRepository
}

func NewPackageService(packageRepo *repositories.PackageRepository, platformRepo *repositories.PlatformRepository) *PackageService {
	return &PackageService{
		packageRepo: packageRepo,
		platformRepo: platformRepo,
	}
}

// CreatePackage creates a new package
func (s *PackageService) CreatePackage(ctx context.Context, pkg *models.Package) error {
	exist, err := s.packageRepo.PackageNameExists(ctx, pkg.Name)
	if err != nil {
		return errors.New("database error")
	}
	if exist {
		return ErrPackageExists
	}

	// Check if the platform exists
	platformExists, err := s.platformRepo.PlatformExists(ctx, pkg.PlatformID)
	if err != nil {
		return errors.New("database error")
	}
	if !platformExists {
		return errors.New("platform does not exist")
	}

	pkg.CreatedAt = time.Now()
	pkg.UpdatedAt = time.Now()

	return s.packageRepo.CreatePackage(ctx, pkg)
}

// GetPackageByID retrieves a package by its ID
func (s *PackageService) GetPackageByID(ctx context.Context, id int) (*models.Package, error) {
	// Check if the package exists in the database
	return s.packageRepo.FindPackageByID(ctx, id)
}

// UpdatePackage updates a package
func (s *PackageService) UpdatePackage(ctx context.Context, pkg *models.Package) error {
	exists, err := s.packageRepo.PackageExists(ctx, pkg.ID)
	if err != nil {
		return errors.New("database error")
	}
	if exists == false {
		return ErrPackageNotFound
	}
	pkg.UpdatedAt = time.Now()
	return s.packageRepo.UpdatePackage(ctx, pkg)
}

// DeletePackage deletes a package
func (s *PackageService) DeletePackage(ctx context.Context, id int) error {
	exists, err := s.packageRepo.PackageExists(ctx, id)
	if err != nil {
		return errors.New("database error")
	}
	if exists == false {
		return ErrPackageNotFound
	}
	return s.packageRepo.DeletePackage(ctx, id)
}

// ListPackagesByPlatform retrieves all packages for a specific platform.
func (s *PackageService) ListPackagesByPlatform(ctx context.Context, platformID int) ([]models.Package, error) {
	return s.packageRepo.ListPackagesByPlatform(ctx, platformID)
}

// PackageType exists in the database
func (s *PackageService) PackageTypeExists(ctx context.Context, pkgType string) (bool, error) {
	return s.packageRepo.PackageTypeExists(ctx, pkgType)
}
