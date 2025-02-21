package services

import (
	"context"
	"errors"
	"time"

	"github.com/mahabub618/minipack/internal/models"
	"github.com/mahabub618/minipack/internal/repositories"
)

var (
	ErrSubscriptionNotFound = errors.New("subscription not found")
)

type SubscriptionService struct {
	subscriptionRepo *repositories.SubscriptionRepository
}

func NewSubscriptionService(subscriptionRepo *repositories.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{subscriptionRepo: subscriptionRepo}
}

// CreateSubscription creates a new subscription.
func (s *SubscriptionService) CreateSubscription(ctx context.Context, sub *models.Subscription) error {
	return s.subscriptionRepo.CreateSubscription(ctx, sub)
}

// GetSubscriptionByID retrieves a subscription by its ID.
func (s *SubscriptionService) GetSubscriptionByID(ctx context.Context, id int) (*models.Subscription, error) {
	return s.subscriptionRepo.FindSubscriptionByID(ctx, id)
}

// UpdateSubscription updates a subscription's details.
func (s *SubscriptionService) UpdateSubscription(ctx context.Context, sub *models.Subscription) error {
	sub.UpdatedAt = time.Now()
	return s.subscriptionRepo.UpdateSubscription(ctx, sub)
}

// DeleteSubscription deletes a subscription by its ID.
func (s *SubscriptionService) DeleteSubscription(ctx context.Context, id int) error {
	return s.subscriptionRepo.DeleteSubscription(ctx, id)
}

// ListSubscriptionsByUser retrieves all subscriptions for a specific user.
func (s *SubscriptionService) ListSubscriptionsByUser(ctx context.Context, userID int) ([]models.Subscription, error) {
	return s.subscriptionRepo.ListSubscriptionsByUser(ctx, userID)
}

// SubscriptionExists checks if a subscription exists in the database.
func (s *SubscriptionService) SubscriptionExists(ctx context.Context, id int) (bool, error) {
	return s.subscriptionRepo.SubscriptionExists(ctx, id)
}
