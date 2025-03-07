package models

import "time"

type Subscription struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id" validate:"required"`
	ValidityID      int       `json:"validity_id" validate:"required"`
	Price           float64   `json:"price"`
	Currency        string    `json:"currency"`
	Status          string    `json:"status"`
	ClientSecret    string    `json:"client_secret"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	AutoRenew       bool      `json:"auto_renew"`
	TrialEndDate    time.Time `json:"trial_end_date"`
	NextBillingDate time.Time `json:"next_billing_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
