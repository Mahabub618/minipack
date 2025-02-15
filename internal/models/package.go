package models

import "time"

type Package struct {
	ID                int       `json:"id,omitempty"`
	PlatformID        int       `json:"platform_id" validate:"required"`
	Name              string    `json:"name" validate:"required"`
	Type              string    `json:"type" validate:"required,oneof=weekly monthly yearly"`
	Duration          int       `json:"duration" validate:"required"`
	Price             float64   `json:"price" validate:"required"`
	Currency          string    `json:"currency" validate:"required"`
	DiscountAmount    float64   `json:"discount_amount"`
	DiscountExpiresAt time.Time `json:"discount_expires_at"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
