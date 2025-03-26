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
		INSERT INTO validities (platform_id, name, description, duration, price, label, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	return repo.db.QueryRow(
		ctx,
		query,
		validity.PlatformID,
		validity.Name,
		validity.Description,
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
		SET platform_id=$1, name=$2, description=$3, duration=$4, price=$5, label=$6, updated_at=$7
		WHERE id = $8`

	_, err := repo.db.Exec(
		ctx,
		query,
		validity.PlatformID,
		validity.Name,
		validity.Description,
		validity.Duration,
		validity.Price,
		validity.Label,
		validity.UpdatedAt,
		validity.ID,
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
		SELECT platform_id, name, description, duration, price, label, created_at, updated_at
		FROM validities
		WHERE id = $1`
	validity := &models.Validity{}
	err := repo.db.QueryRow(ctx, query, id).Scan(
		&validity.PlatformID,
		&validity.Name,
		&validity.Description,
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
		SELECT id, platform_id, name, description, duration, price, label,  created_at, updated_at
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
			&validity.Name,
			&validity.Description,
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

func (repo *ValidityRepository) ValidityLabelExists(ctx context.Context, label string, platformId int) (bool, error) {
	query := `
		SELECT EXISTS(SELECT id FROM validities WHERE label = $1 AND platform_id = $2)`
	var exists bool
	err := repo.db.QueryRow(ctx, query, label, platformId).Scan(&exists)
	return exists, err
}
