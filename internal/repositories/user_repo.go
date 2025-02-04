package repositories

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mahabub618/minipack/internal/models"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser inserts a new user into the database.
func (repo *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (name, email, password, phone, address_street, address_city, address_country, role, is_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`

	row := repo.db.QueryRow(ctx, query, user.Name, user.Email, user.Password, user.Phone, user.AddressStreet, user.AddressCity, user.AddressCountry, user.Role, user.IsVerified, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
	return row
}

// FindUserByEmail retrieves a user by email
func (repo *UserRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, name, email, password, phone, address_street, address_city, address_country, role, is_verified, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	user := &models.User{}
	row := repo.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone, &user.AddressStreet, &user.AddressCity, &user.AddressCountry, &user.Role, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt)
	if row == pgx.ErrNoRows {
		return nil, nil
	}
	return user, nil
}

// EmailExists check if an email already exists in the database
func (repo *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	query := `
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)
	`
	var exists bool
	err := repo.db.QueryRow(ctx, query, email).Scan(&exists)
	return exists, err
}
