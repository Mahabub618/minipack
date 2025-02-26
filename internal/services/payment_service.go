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
	ErrPaymentNotFound = errors.New("payment not found")
)

type PaymentService struct {
	paymentRepo *repositories.PaymentRepository
}

func NewPaymentService(paymentRepo *repositories.PaymentRepository) *PaymentService {
	return &PaymentService{paymentRepo: paymentRepo}
}

// CreatePayment creates a new payment using Stripe.
func (s *PaymentService) CreatePayment(ctx context.Context, payment *models.Payment) error {
	// Set the Stripe API key from the environment
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	// Create a PaymentIntent with Stripe
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(payment.Amount * 100)),
		Currency: stripe.String(payment.Currency),
		PaymentMethodTypes: []*string{
			stripe.String("card"),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return err
	}

	// Set the transaction ID from Stripe
	payment.TransactionID = pi.ID
	payment.Status = "pending"
	payment.CreatedAt = time.Now()

	return s.paymentRepo.CreatePayment(ctx, payment)
}

// ConfirmPayment confirms a payment using Stripe.
func (s *PaymentService) ConfirmPayment(ctx context.Context, paymentID int) error {
	payment, err := s.paymentRepo.FindPaymentByID(ctx, paymentID)
	if err != nil {
		return err
	}
	if payment == nil {
		return ErrPaymentNotFound
	}

	// Confirm the PaymentIntent with Stripe
	pi, err := paymentintent.Confirm(
		payment.TransactionID,
		nil,
	)
	if err != nil {
		return err
	}

	// Update the payment status based on Stripe's response
	payment.Status = string(pi.Status)
	return s.paymentRepo.UpdatePayment(ctx, payment)
}

// GetPaymentByID retrieves a payment by its ID.
func (s *PaymentService) GetPaymentByID(ctx context.Context, id int) (*models.Payment, error) {
	return s.paymentRepo.FindPaymentByID(ctx, id)
}

// ListPaymentsBySubscription retrieves all payments for a specific subscription.
func (s *PaymentService) ListPaymentsBySubscription(ctx context.Context, subscriptionID int) ([]models.Payment, error) {
	return s.paymentRepo.ListPaymentsBySubscription(ctx, subscriptionID)
}
