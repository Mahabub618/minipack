package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int            `db:"id"`
	Name           string         `db:"name"`
	Email          string         `db:"email"`
	Password       string         `db:"password"`
	Phone          sql.NullString `db:"phone"`
	AddressStreet  sql.NullString `db:"address_street"`
	AddressCity    sql.NullString `db:"address_city"`
	AddressCountry sql.NullString `db:"address_country"`
	Role           string         `db:"role"`
	IsVerified     bool           `db:"is_verified"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at"`
}

type RefreshToken struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
}
