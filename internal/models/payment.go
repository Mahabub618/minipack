package models

import "time"

type Payment struct {
	ID             int       `json:"id"`
	SubscriptionID int       `json:"subscription_id" validate:"required"`
	Amount         float64   `json:"amount" validate:"required"`
	Currency       string    `json:"currency" validate:"required"`
	PaymentMethod  string    `json:"payment_method" validate:"required"`
	TransactionID  string    `json:"transaction_id" validate:"required"`
	Status         string    `json:"status" validate:"oneof=success failed pending refunded"`
	CreatedAt      time.Time `json:"created_at"`
}
