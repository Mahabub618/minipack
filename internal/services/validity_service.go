package services

import (
	"context"
	"github.com/mahabub618/minipack/internal/models"
	"github.com/mahabub618/minipack/internal/repositories"
	"time"
)

type ValidityService struct {
	validityRepository *repositories.ValidityRepository
}

func NewValidityService(validityRepository *repositories.ValidityRepository) *ValidityService {
	return &ValidityService{validityRepository: validityRepository}
}

func (s *ValidityService) CreateValidity(ctx context.Context, validity *models.Validity) error {
	now := time.Now()
	validity.CreatedAt, validity.UpdatedAt = now, now
	return s.validityRepository.CreateValidityRepository(ctx, validity)
}

func (s *ValidityService) UpdateValidity(ctx context.Context, validity *models.Validity) error {
	validity.UpdatedAt = time.Now()
	return s.validityRepository.UpdateValidityRepository(ctx, validity)
}

func (s *ValidityService) GetValidityByID(ctx context.Context, id int) (*models.Validity, error) {
	return s.validityRepository.GetValidityByID(ctx, id)
}

func (s *ValidityService) DeleteValidityByID(ctx context.Context, id int) error {
	return s.validityRepository.DeleteValidityRepository(ctx, id)
}

func (s *ValidityService) GetAllValiditiesByPlatformID(ctx context.Context, id int) ([]*models.Validity, error) {
	return s.validityRepository.ListValiditiesByPlatformID(ctx, id)
}
