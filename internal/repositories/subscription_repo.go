package repositories

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mahabub618/minipack/internal/models"
)

type SubscriptionRepository struct {
	db *pgxpool.Pool
}

func NewSubscriptionRepository(db *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

// CreateSubscription inserts a new subscription into the database.
func (r *SubscriptionRepository) CreateSubscription(ctx context.Context, sub *models.Subscription) error {
	query := `
		INSERT INTO subscriptions (user_id, package_id, price, currency, status, client_secret, start_date, end_date, auto_renew, trial_end_date, next_billing_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id
	`
	return r.db.QueryRow(
		ctx,
		query,
		sub.UserID,
		sub.PackageID,
		sub.Price,
		sub.Currency,
		sub.Status,
		sub.ClientSecret,
		sub.StartDate,
		sub.EndDate,
		sub.AutoRenew,
		sub.TrialEndDate,
		sub.NextBillingDate,
		sub.CreatedAt,
		sub.UpdatedAt,
	).Scan(&sub.ID)
}

// FindSubscriptionByID retrieves a subscription by its ID.
func (r *SubscriptionRepository) FindSubscriptionByID(ctx context.Context, id int) (*models.Subscription, error) {
	query := `
		SELECT id, user_id, package_id, start_date, end_date, status, auto_renew, trial_end_date, next_billing_date, created_at, updated_at
		FROM subscriptions
		WHERE id = $1
	`
	sub := &models.Subscription{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&sub.ID,
		&sub.UserID,
		&sub.PackageID,
		&sub.StartDate,
		&sub.EndDate,
		&sub.Status,
		&sub.AutoRenew,
		&sub.TrialEndDate,
		&sub.NextBillingDate,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return sub, nil
}

// UpdateSubscription updates a subscription's details.
func (r *SubscriptionRepository) UpdateSubscription(ctx context.Context, sub *models.Subscription) error {
	// Start building the query
	query := `
		UPDATE subscriptions
		SET user_id = $1, package_id = $2, start_date = $3, end_date = $4, status = $5, auto_renew = $6, trial_end_date = $7, next_billing_date = $8, updated_at = $9
		WHERE id = $10
	`
	_, err := r.db.Exec(
		ctx,
		query,
		sub.UserID,
		sub.PackageID,
		sub.StartDate,
		sub.EndDate,
		sub.Status,
		sub.AutoRenew,
		sub.TrialEndDate,
		sub.NextBillingDate,
		sub.UpdatedAt,
		sub.ID,
	)
	return err
}

// DeleteSubscription deletes a subscription by its ID.
func (r *SubscriptionRepository) DeleteSubscription(ctx context.Context, id int) error {
	query := `
		DELETE FROM subscriptions
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// ListSubscriptionsByUser retrieves all subscriptions for a specific user.
func (r *SubscriptionRepository) ListSubscriptionsByUser(ctx context.Context, userID int) ([]models.Subscription, error) {
	query := `
		SELECT id, user_id, package_id, start_date, end_date, status, auto_renew, trial_end_date, next_billing_date, created_at, updated_at
		FROM subscriptions
		WHERE user_id = $1
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []models.Subscription
	for rows.Next() {
		var sub models.Subscription
		if err := rows.Scan(
			&sub.ID,
			&sub.UserID,
			&sub.PackageID,
			&sub.StartDate,
			&sub.EndDate,
			&sub.Status,
			&sub.AutoRenew,
			&sub.TrialEndDate,
			&sub.NextBillingDate,
			&sub.CreatedAt,
			&sub.UpdatedAt,
		); err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, sub)
	}
	return subscriptions, nil
}

// SubscriptionExists checks if a subscription exists in the database.
func (r *SubscriptionRepository) SubscriptionExists(ctx context.Context, id int) (bool, error) {
	query := `
		SELECT exists(SELECT 1 FROM subscriptions WHERE id = $1)
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
