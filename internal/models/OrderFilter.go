package models

import (
	"time"
)

type OrderFilter struct {
	UserID        *string
	Status        *string
	PaymentStatus *string
	StartDate     *time.Time
	EndDate       *time.Time
	MinTotal      *float64
	MaxTotal      *float64
	Limit         int
	Offset        int
}
