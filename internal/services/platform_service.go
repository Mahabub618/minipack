package services

import (
	"context"
	"errors"
	"github.com/mahabub618/minipack/internal/models"
	"github.com/mahabub618/minipack/internal/repositories"
	"time"
)

var ErrPlatformExists = errors.New("platform name already exists")
var ErrPlatformNotFound = errors.New("platform not found")

type PlatformService struct {
	platformRepo *repositories.PlatformRepository
}

func NewPlatformService(platformRepo *repositories.PlatformRepository) *PlatformService {
	return &PlatformService{platformRepo: platformRepo}
}

// CreatePlatform creates a new platform
func (s *PlatformService) CreatePlatform(ctx context.Context, platform *models.Platform) error {
	exist, err := s.platformRepo.PlatformNameExists(ctx, platform.Name)
	if err != nil {
		return errors.New("database error")
	}
	if exist {
		return ErrPlatformExists
	}

	platform.Status = "inactive"
	platform.CreatedAt = time.Now()
	platform.UpdatedAt = time.Now()

	return s.platformRepo.CreatePlatform(ctx, platform)
}

// GetPlatformByID retrieves a platform by its ID
func (s *PlatformService) GetPlatformByID(ctx context.Context, id int) (*models.Platform, error) {
	return s.platformRepo.FindPlatformByID(ctx, id)
}

// UpdatePlatform updates a platform
func (s *PlatformService) UpdatePlatform(ctx context.Context, platform *models.Platform) error {
	exists, err := s.platformRepo.FindPlatformByID(ctx, platform.ID)
	if err != nil {
		return errors.New("database error")
	}
	if exists == nil {
		return ErrPlatformNotFound
	}
	platform.UpdatedAt = time.Now()
	return s.platformRepo.UpdatePlatform(ctx, platform)
}

// DeletePlatform deletes a platform
func (s *PlatformService) DeletePlatform(ctx context.Context, id int) error {
	exists, err := s.platformRepo.FindPlatformByID(ctx, id)
	if err != nil {
		return errors.New("database error")
	}
	if exists == nil {
		return ErrPlatformNotFound
	}
	return s.platformRepo.DeletePlatform(ctx, id)
}

// ActivatePlatform activates a platform
func (s *PlatformService) ActivatePlatform(ctx context.Context, id int) error {
	exists, err := s.platformRepo.FindPlatformByID(ctx, id)
	if err != nil {
		return errors.New("database error")
	}
	if exists == nil {
		return ErrPlatformNotFound
	}
	platform := &models.Platform{
		ID:     id,
		Status: "active",
	}
	return s.platformRepo.UpdatePlatform(ctx, platform)
}

// DeactivatePlatform deactivates a platform
func (s *PlatformService) DeactivatePlatform(ctx context.Context, id int) error {
	exists, err := s.platformRepo.FindPlatformByID(ctx, id)
	if err != nil {
		return errors.New("database error")
	}
	if exists == nil {
		return ErrPlatformNotFound
	}
	platform := &models.Platform{
		ID:     id,
		Status: "inactive",
	}
	return s.platformRepo.UpdatePlatform(ctx, platform)
}

// ListPlatforms retrieves a list of platforms
func (s *PlatformService) ListPlatforms(ctx context.Context) ([]models.Platform, error) {
	return s.platformRepo.ListPlatforms(ctx)
}
