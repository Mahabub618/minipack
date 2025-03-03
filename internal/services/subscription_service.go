package services

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/mahabub618/minipack/internal/models"
	"github.com/mahabub618/minipack/internal/repositories"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
)

var (
	ErrSubscriptionNotFound = errors.New("subscription not found")
)

type SubscriptionService struct {
	subscriptionRepo *repositories.SubscriptionRepository
	packageRepo      *repositories.PackageRepository
}

func NewSubscriptionService(subscriptionRepo *repositories.SubscriptionRepository, packageRepo *repositories.PackageRepository) *SubscriptionService {
	return &SubscriptionService{
		subscriptionRepo: subscriptionRepo,
		packageRepo:      packageRepo,
	}
}

// CreateSubscription creates a new subscription.
func (s *SubscriptionService) CreateSubscription(ctx context.Context, sub *models.Subscription) error {
	// Set the Stripe API key from the environment
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Fetch package details
	pkg, err := s.packageRepo.FindPackageByID(ctx, sub.PackageID)
	if err != nil {
		return errors.New("package not found")
	}

	// Create a PaymentIntent with Stripe
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(pkg.Price * 100)),
		Currency: stripe.String(pkg.Currency),
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		var stripeErr *stripe.Error
		if errors.As(err, &stripeErr) {
			return errors.New("A stripe error occurred: " + stripeErr.Error())
		}
	}
	// Set subscription details
	now := time.Now()
	sub.CreatedAt, sub.UpdatedAt, sub.StartDate = now, now, now
	sub.EndDate = now.Add(time.Duration(pkg.Duration) * 24 * time.Hour)
	sub.Price, sub.Currency = pkg.Price-pkg.DiscountAmount, pkg.Currency
	sub.Status = "pending"
	sub.ClientSecret = pi.ClientSecret

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
