package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mahabub618/minipack/internal/models"
)

type PlatformRepository struct {
	db *pgxpool.Pool
}

func NewPlatformRepository(db *pgxpool.Pool) *PlatformRepository {
	return &PlatformRepository{db: db}
}

// CreatePlatform creates a new platform into the database
func (repo *PlatformRepository) CreatePlatform(ctx context.Context, platform *models.Platform) error {
	query := `
        INSERT INTO platforms (name, description, logo_url, status, supported_countries, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id`

	row := repo.db.QueryRow(ctx, query, platform.Name, platform.Description, platform.LogoURL, platform.Status, platform.SupportedCountries, platform.CreatedAt, platform.UpdatedAt).Scan(&platform.ID)

	return row
}

// FindPlatformByID retrieves a platform by its ID.
func (repo *PlatformRepository) FindPlatformByID(ctx context.Context, id int) (*models.Platform, error) {
	query := `
        SELECT id, name, description, logo_url, status, supported_countries, created_at, updated_at
        FROM platforms
        WHERE id = $1
    `
	platform := &models.Platform{}
	row := repo.db.QueryRow(ctx, query, id).Scan(
		&platform.ID,
		&platform.Name,
		&platform.Description,
		&platform.LogoURL,
		&platform.Status,
		&platform.SupportedCountries,
		&platform.CreatedAt,
		&platform.UpdatedAt)
	if row == pgx.ErrNoRows {
		return nil, nil
	}
	return platform, nil
}

// UpdatePlatform updates a platform in the database
func (repo *PlatformRepository) UpdatePlatform(ctx context.Context, platform *models.Platform) error {
	query := `
        UPDATE platforms
        SET name = $1, description = $2, logo_url = $3, status = $4, supported_countries = $5, updated_at = $6
        WHERE id = $7`

	_, err := repo.db.Exec(ctx, query, platform.Name, platform.Description, platform.LogoURL, platform.Status, platform.SupportedCountries, platform.UpdatedAt, platform.ID)
	return err
}

// DeletePlatform deletes a platform from the database
func (repo *PlatformRepository) DeletePlatform(ctx context.Context, id int) error {
	query := `
        DELETE FROM platforms
        WHERE id = $1
    `
	_, err := repo.db.Exec(ctx, query, id)
	return err
}

// ListPlatforms retrieves a list of platforms from the database
func (r *PlatformRepository) ListPlatforms(ctx context.Context) ([]models.Platform, error) {
	query := `
		SELECT id, name, description, logo_url, status, supported_countries, created_at, updated_at
		FROM platforms
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var platforms []models.Platform
	for rows.Next() {
		var platform models.Platform
		if err := rows.Scan(
			&platform.ID,
			&platform.Name,
			&platform.Description,
			&platform.LogoURL,
			&platform.Status,
			&platform.SupportedCountries,
			&platform.CreatedAt,
			&platform.UpdatedAt); err != nil {
			return nil, err
		}
		platforms = append(platforms, platform)
	}
	return platforms, nil
}

func (r *PlatformRepository) ListPlatformsWithPriceRanges(ctx context.Context) ([]models.Platform, error) {
	query := `
        SELECT 
            p.id, 
            p.name, 
            p.description, 
            p.logo_url, 
            p.status, 
            p.supported_countries, 
            p.created_at, 
            p.updated_at,
            MIN(v.price) as min_price,
            MAX(v.price) as max_price
        FROM platforms p
        LEFT JOIN validities v ON p.id = v.platform_id
        GROUP BY p.id, p.name, p.description, p.logo_url, p.status, p.supported_countries, p.created_at, p.updated_at
    `

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var platforms []models.Platform
	for rows.Next() {
		var platform models.Platform
		var minPrice, maxPrice sql.NullFloat64

		if err := rows.Scan(
			&platform.ID,
			&platform.Name,
			&platform.Description,
			&platform.LogoURL,
			&platform.Status,
			&platform.SupportedCountries,
			&platform.CreatedAt,
			&platform.UpdatedAt,
			&minPrice,
			&maxPrice,
		); err != nil {
			return nil, err
		}

		// Format price range if available
		if minPrice.Valid && maxPrice.Valid {
			if minPrice.Float64 == maxPrice.Float64 {
				platform.PriceRange = fmt.Sprintf("%.2f ৳", minPrice.Float64)
			} else {
				platform.PriceRange = fmt.Sprintf("%.2f ৳ - %.2f ৳", minPrice.Float64, maxPrice.Float64)
			}
		}

		platforms = append(platforms, platform)
	}
	return platforms, nil
}

// PlatformExists checks if a platform exists in the database
func (repo *PlatformRepository) PlatformExists(ctx context.Context, id int) (bool, error) {
	query := `
        SELECT EXISTS(SELECT 1 FROM platforms WHERE id = $1)
    `
	var exists bool
	err := repo.db.QueryRow(ctx, query, id).Scan(&exists)
	return exists, err
}

// PlatformNameExists checks if a platform name already exists in the database
func (repo *PlatformRepository) PlatformNameExists(ctx context.Context, name string) (bool, error) {
	query := `
        SELECT EXISTS(SELECT 1 FROM platforms WHERE name = $1)
    `
	var exists bool
	err := repo.db.QueryRow(ctx, query, name).Scan(&exists)
	return exists, err
}
