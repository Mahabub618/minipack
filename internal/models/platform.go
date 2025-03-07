package models

import "time"

type Platform struct {
	ID                 int       `json:"id"`
	Name               string    `json:"name" validate:"required"`
	Description        string    `json:"description"`
	LogoURL            string    `json:"logo_url"`
	Status             string    `json:"status" validate:"oneof=active inactive"`
	SupportedCountries []string  `json:"supported_countries"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
