package models

import "time"

type Validity struct {
	ID          int       `json:"id,omitempty"`
	PlatformID  int       `json:"platform_id" validate:"required"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Duration    int       `json:"duration" validate:"required"`
	Price       float64   `json:"price" validate:"required"`
	Label       string    `json:"label" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
