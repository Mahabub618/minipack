package models

import "time"

type Subscription struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id" validate:"required"`
	PackageID       int       `json:"package_id" validate:"required"`
	Price           float64   `json:"price" validate:"required"`
	Currency        string    `json:"currency" validate:"required"`
	StartDate       time.Time `json:"start_date validate:"required"`
	EndDate         time.Time `json:"end_date validate:"required"`
	Status          string    `json:"status" validate:"oneof=active expired cancelled paused"`
	AutoRenew       bool      `json:"auto_renew"`
	TrialEndDate    time.Time `json:"trial_end_date"`
	NextBillingDate time.Time `json:"next_billing_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
