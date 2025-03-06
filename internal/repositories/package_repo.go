package repositories

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mahabub618/minipack/internal/models"
)

type PackageRepository struct {
	db *pgxpool.Pool
}

func NewPackageRepository(db *pgxpool.Pool) *PackageRepository {
	return &PackageRepository{db: db}
}

// FindPackageByID retrieves a package by its ID.
func (repo *PackageRepository) FindPackageByID(ctx context.Context, id int) (*models.Package, error) {
	query := `
		SELECT id, platform_id, discount_amount, discount_expires_at, created_at, updated_at
		FROM packages
		WHERE id = $1
	`
	pkg := &models.Package{}
	row := repo.db.QueryRow(ctx, query, id).Scan(
		&pkg.ID,
		&pkg.PlatformID,
		&pkg.DiscountAmount,
		&pkg.DiscountExpiresAt,
		&pkg.CreatedAt,
		&pkg.UpdatedAt)
	if row == pgx.ErrNoRows {
		return nil, nil
	}
	return pkg, nil
}

// ListPackagesByPlatform retrieves all packages for a specific platform.
func (r *PackageRepository) ListPackagesByPlatform(ctx context.Context, platformID int) ([]models.Package, error) {
	query := `
		SELECT id, platform_id, discount_amount, discount_expires_at, created_at, updated_at
		FROM packages
		WHERE platform_id = $1
	`
	rows, err := r.db.Query(ctx, query, platformID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var packages []models.Package
	for rows.Next() {
		var pkg models.Package
		if err := rows.Scan(
			&pkg.ID,
			&pkg.PlatformID,
			&pkg.DiscountAmount,
			&pkg.DiscountExpiresAt,
			&pkg.CreatedAt,
			&pkg.UpdatedAt,
		); err != nil {
			return nil, err
		}
		packages = append(packages, pkg)
	}
	return packages, nil
}
