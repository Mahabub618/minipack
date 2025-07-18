package repositories

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mahabub618/minipack/internal/models"
	"time"
)

type AnonymousUserRepository struct {
	db *pgxpool.Pool
}

func NewAnonymousUserRepository(db *pgxpool.Pool) *AnonymousUserRepository {
	return &AnonymousUserRepository{db: db}
}

func (r *AnonymousUserRepository) Create(ctx context.Context, sessionID string) (*models.AnonymousUser, error) {
	user := &models.AnonymousUser{
		ID:        uuid.New().String(),
		SessionID: sessionID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(180 * 24 * time.Hour), // 180 days expiration
	}

	query := `INSERT INTO anonymous_users (id, session_id, created_at, expires_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, query, user.ID, user.SessionID, user.CreatedAt, user.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *AnonymousUserRepository) GetBySessionID(ctx context.Context, sessionID string) (*models.AnonymousUser, error) {
	query := `SELECT id, session_id, created_at, expires_at FROM anonymous_users WHERE session_id = $1`
	row := r.db.QueryRow(ctx, query, sessionID)

	var user models.AnonymousUser
	err := row.Scan(&user.ID, &user.SessionID, &user.CreatedAt, &user.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *AnonymousUserRepository) MigrateToUser(ctx context.Context, anonymousUserID, registeredUserID string) error {
	// Update all carts and other references
	_, err := r.db.Exec(ctx,
		`UPDATE carts SET user_id = $1 WHERE user_id = $2`,
		registeredUserID, anonymousUserID)
	return err
}
