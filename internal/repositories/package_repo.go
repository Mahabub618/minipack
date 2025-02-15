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

// CreatePackage creates a new package into the database
func (repo *PackageRepository) CreatePackage(ctx context.Context, pkg *models.Package) error {
	query := `
		INSERT INTO packages (platform_id, name, type, duration, price, currency, discount_amount, discount_expires_at, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`

	row := repo.db.QueryRow(
		ctx,
		query,
		pkg.PlatformID,
		pkg.Name,
		pkg.Type,
		pkg.Duration,
		pkg.Price,
		pkg.Currency,
		pkg.DiscountAmount,
		pkg.DiscountExpiresAt,
		pkg.CreatedAt,
		pkg.UpdatedAt).Scan(&pkg.ID)

	return row
}

// FindPackageByID retrieves a package by its ID.
func (repo *PackageRepository) FindPackageByID(ctx context.Context, id int) (*models.Package, error) {
	query := `
		SELECT id, platform_id, name, type, duration, price, currency, discount_amount, discount_expires_at, created_at, updated_at
		FROM packages
		WHERE id = $1
	`
	pkg := &models.Package{}
	row := repo.db.QueryRow(ctx, query, id).Scan(
		&pkg.ID,
		&pkg.PlatformID,
		&pkg.Name,
		&pkg.Type,
		&pkg.Duration,
		&pkg.Price,
		&pkg.Currency,
		&pkg.DiscountAmount,
		&pkg.DiscountExpiresAt,
		&pkg.CreatedAt,
		&pkg.UpdatedAt)
	if row == pgx.ErrNoRows {
		return nil, nil
	}
	return pkg, nil
}

// UpdatePackage updates a package in the database
func (repo *PackageRepository) UpdatePackage(ctx context.Context, pkg *models.Package) error {
	query := `
		UPDATE packages
		SET platform_id = $1, name = $2, type = $3, duration = $4, price = $5, currency = $6, discount_amount = $7, discount_expires_at = $8, updated_at = $9
		WHERE id = $10`

	_, err := repo.db.Exec(ctx, query, pkg.PlatformID, pkg.Name, pkg.Type, pkg.Duration, pkg.Price, pkg.Currency, pkg.DiscountAmount, pkg.DiscountExpiresAt, pkg.UpdatedAt, pkg.ID)
	return err
}

// DeletePackage deletes a package from the database
func (repo *PackageRepository) DeletePackage(ctx context.Context, id int) error {
	query := `
		DELETE FROM packages
		WHERE id = $1
	`
	_, err := repo.db.Exec(ctx, query, id)
	return err
}

// ListPackagesByPlatform retrieves all packages for a specific platform.
func (r *PackageRepository) ListPackagesByPlatform(ctx context.Context, platformID int) ([]models.Package, error) {
	query := `
		SELECT id, platform_id, name, type, duration, price, currency, discount_amount, discount_expires_at, created_at, updated_at
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
			&pkg.Name,
			&pkg.Type,
			&pkg.Duration,
			&pkg.Price,
			&pkg.Currency,
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

// PackageExists checks if a package exists in the database
func (repo *PackageRepository) PackageExists(ctx context.Context, id int) (bool, error) {
	query := `
		SELECT exists(SELECT 1 FROM packages WHERE id = $1)
	`
	var exists bool
	err := repo.db.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// PackageNameExists checks if a package name already exists in the database
func (repo *PackageRepository) PackageNameExists(ctx context.Context, name string) (bool, error) {
	query := `
		SELECT exists(SELECT 1 FROM packages WHERE name = $1)
	`
	var exists bool
	err := repo.db.QueryRow(ctx, query, name).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// PackageTypeExists checks if a package type already exists in the database
func (repo *PackageRepository) PackageTypeExists(ctx context.Context, pkgType string) (bool, error) {
	query := `
		SELECT exists(SELECT 1 FROM packages WHERE type = $1)
	`
	var exists bool
	err := repo.db.QueryRow(ctx, query, pkgType).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
