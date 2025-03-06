package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mahabub618/minipack/internal/models"
)

type ValidityRepository struct {
	db *pgxpool.Pool
}

func NewValidityRepository(db *pgxpool.Pool) *ValidityRepository { return &ValidityRepository{db: db} }

// CreateValidityRepository creates a new validity into the database for a specific platform
func (repo *ValidityRepository) CreateValidityRepository(ctx context.Context, validity *models.Validity) error {
	query := `
		INSERT INTO validities (platform_id, duration, price, label, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	return repo.db.QueryRow(
		ctx,
		query,
		validity.PlatformID,
		validity.Duration,
		validity.Price,
		validity.Label,
		validity.CreatedAt,
		validity.UpdatedAt,
	).Scan(&validity.ID, &validity.CreatedAt, &validity.UpdatedAt)
}

func (repo *ValidityRepository) UpdateValidityRepository(ctx context.Context, validity *models.Validity) error {
	query := `
		UPDATE validities
		SET platform_id=$1, duration=$2, price=$3, label=$4, created_at=$5, updated_at=$6
		WHERE id = $7`

	_, err := repo.db.Exec(
		ctx,
		query,
		validity.PlatformID,
		validity.Duration,
		validity.Price,
		validity.Label,
		validity.CreatedAt,
		validity.UpdatedAt,
	)
	return err
}

func (repo *ValidityRepository) DeleteValidityRepository(ctx context.Context, id int) error {
	query := `DELETE FROM validities WHERE id = $1`
	_, err := repo.db.Exec(ctx, query, id)
	return err
}

func (repo *ValidityRepository) GetValidityByID(ctx context.Context, id int) (*models.Validity, error) {
	query := `
		SELECT platform_id, duration, price, label, created_at, updated_at
		FROM validities
		WHERE id = $1`
	validity := &models.Validity{}
	err := repo.db.QueryRow(ctx, query, id).Scan(
		&validity.PlatformID,
		&validity.Duration,
		&validity.Price,
		&validity.Label,
		&validity.CreatedAt,
		&validity.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return validity, nil
}

// List all ValidityByPlatformID
func (repo *ValidityRepository) ListValiditiesByPlatformID(ctx context.Context, platformID int) ([]*models.Validity, error) {
	query := `
		SELECT id, platform_id, duration, price, label,  created_at, updated_at
		FROM validities
		WHERE platform_id = $1`
	rows, err := repo.db.Query(ctx, query, platformID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var validities []*models.Validity
	for rows.Next() {
		validity := &models.Validity{}
		if err := rows.Scan(
			&validity.ID,
			&validity.PlatformID,
			&validity.Duration,
			&validity.Price,
			&validity.Label,
			&validity.CreatedAt,
			&validity.UpdatedAt,
		); err != nil {
			return nil, err
		}
		validities = append(validities, validity)
	}
	return validities, nil
}

func (repo *ValidityRepository) ValidityLabelExists(ctx context.Context, label string) (bool, error) {
	query := `
		SELECT EXISTS(SELECT id FROM validities WHERE label = $1)`
	var exists bool
	err := repo.db.QueryRow(ctx, query, label).Scan(&exists)
	return exists, err
}
