package repositories

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mahabub618/minipack/internal/models"
)

type PaymentRepository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// CreatePayment creates a new payment into the database
func (repo *PaymentRepository) CreatePayment(ctx context.Context, payment *models.Payment) error {
	query := `
		INSERT INTO payments (subscription_id, amount, currency, payment_method, transaction_id, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	return repo.db.QueryRow(
		ctx,
		query,
		payment.SubscriptionID,
		payment.Amount,
		payment.Currency,
		payment.PaymentMethod,
		payment.TransactionID,
		payment.Status,
		payment.CreatedAt,
	).Scan(&payment.ID)
}

// FindPaymentByID retrieves a payment by its ID.
func (repo *PaymentRepository) FindPaymentByID(ctx context.Context, id int) (*models.Payment, error) {
	query := `
		SELECT id, subscription_id, amount, currency, payment_method, transaction_id, status, created_at
		FROM payments
		WHERE id = $1`

	payment := &models.Payment{}
	err := repo.db.QueryRow(ctx, query, id).Scan(
		&payment.ID,
		&payment.SubscriptionID,
		&payment.Amount,
		&payment.Currency,
		&payment.PaymentMethod,
		&payment.TransactionID,
		&payment.Status,
		&payment.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

// UpdatePayment updates a payment's details.
func (r *PaymentRepository) UpdatePayment(ctx context.Context, payment *models.Payment) error {
	query := `
		UPDATE payments
		SET subscription_id = $1, amount = $2, currency = $3, payment_method = $4, transaction_id = $5, status = $6
		WHERE id = $7
	`
	_, err := r.db.Exec(
		ctx,
		query,
		payment.SubscriptionID,
		payment.Amount,
		payment.Currency,
		payment.PaymentMethod,
		payment.TransactionID,
		payment.Status,
		payment.ID,
	)
	return err
}

// ListPaymentsBySubscription retrieves all payments for a specific subscription.
func (r *PaymentRepository) ListPaymentsBySubscription(ctx context.Context, subscriptionID int) ([]models.Payment, error) {
	query := `
		SELECT id, subscription_id, amount, currency, payment_method, transaction_id, status, created_at
		FROM payments
		WHERE subscription_id = $1
	`
	rows, err := r.db.Query(ctx, query, subscriptionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		if err := rows.Scan(
			&payment.ID,
			&payment.SubscriptionID,
			&payment.Amount,
			&payment.Currency,
			&payment.PaymentMethod,
			&payment.TransactionID,
			&payment.Status,
			&payment.CreatedAt,
		); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return payments, nil
}
