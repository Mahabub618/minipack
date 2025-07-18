package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mahabub618/minipack/internal/models"
	"log"
	"time"
)

type CartRepository struct {
	db *pgxpool.Pool
}

func NewCartRepository(db *pgxpool.Pool) *CartRepository {
	return &CartRepository{db: db}
}

// CreateCart creates a new cart in the database
func (r *CartRepository) CreateCart(ctx context.Context, cart *models.Cart) error {
	cart.ID = uuid.New().String()
	cart.CreatedAt = time.Now()
	cart.UpdatedAt = time.Now()

	query := `
		INSERT INTO carts (id, user_id, status, created_at, updated_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query,
		cart.ID, cart.UserID, cart.Status, cart.CreatedAt, cart.UpdatedAt, cart.ExpiresAt,
	)
	return err
}

func (r *CartRepository) GetCartByID(ctx context.Context, id string) (*models.Cart, error) {
	query := `
		SELECT id, user_id, status, created_at, updated_at, expires_at
		FROM carts WHERE id = $1
	`
	cart := &models.Cart{}
	row := r.db.QueryRow(ctx, query, id).Scan(
		&cart.ID,
		&cart.UserID,
		&cart.Status,
		&cart.CreatedAt,
		&cart.UpdatedAt,
		&cart.ExpiresAt)
	if row == pgx.ErrNoRows {
		return nil, nil
	}

	return cart, nil
}

func (r *CartRepository) UpdateCart(ctx context.Context, cart *models.Cart) error {
	cart.UpdatedAt = time.Now()
	query := `
		UPDATE carts SET
			status = $2,
			updated_at = $3,
			expires_at = $4
		WHERE id = $1
	`
	_, err := r.db.Exec(
		ctx,
		query,
		cart.ID,
		cart.Status,
		cart.UpdatedAt,
		cart.ExpiresAt)
	return err
}

func (r *CartRepository) DeleteCart(ctx context.Context, id string) error {
	query := `DELETE FROM carts WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *CartRepository) GetCartItems(ctx context.Context, cartID string) ([]models.CartItem, error) {
	query := `
		SELECT id, cart_id, product_id, product_name, quantity, unit_price, discount, total_price, created_at, updated_at
		FROM cart_items WHERE cart_id = $1
	`
	rows, err := r.db.Query(ctx, query, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.CartItem
	for rows.Next() {
		item := models.CartItem{}
		err := rows.Scan(
			&item.ID,
			&item.CartID,
			&item.ProductID,
			&item.ProductName,
			&item.Quantity,
			&item.UnitPrice,
			&item.Discount,
			&item.TotalPrice,
			&item.CreatedAt,
			&item.UpdatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *CartRepository) AddItemToCart(ctx context.Context, item *models.CartItem) error {
	item.ID = uuid.New().String()
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	log.Println("#255 item.ID", item.ID)

	query := `
		INSERT INTO cart_items (
			id, cart_id, product_id, product_name, quantity, unit_price, discount, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(
		ctx,
		query,
		item.ID,
		item.CartID,
		item.ProductID,
		item.ProductName,
		item.Quantity,
		item.UnitPrice,
		item.Discount,
		item.CreatedAt,
		item.UpdatedAt)
	return err
}

func (r *CartRepository) RemoveItemFromCart(ctx context.Context, itemID string) error {
	query := `DELETE FROM cart_items WHERE id = $1`
	_, err := r.db.Exec(ctx, query, itemID)
	return err
}

func (r *CartRepository) UpdateCartItemQuantity(ctx context.Context, itemID string, quantity int) error {
	query := `UPDATE cart_items SET quantity = $2, updated_at = $3 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, itemID, quantity, time.Now())
	return err
}

func (r *CartRepository) GetActiveCart(ctx context.Context, userID string) (*models.Cart, error) {
	query := `
        SELECT id, user_id, status, created_at, updated_at, expires_at
        FROM carts 
        WHERE user_id = $1 AND status = 'active'
        ORDER BY created_at DESC
        LIMIT 1
    `
	row := r.db.QueryRow(ctx, query, userID)

	var cart models.Cart
	err := row.Scan(
		&cart.ID, &cart.UserID, &cart.Status, &cart.CreatedAt, &cart.UpdatedAt, &cart.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &cart, nil
}
