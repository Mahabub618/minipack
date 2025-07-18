package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mahabub618/minipack/internal/models"
	"time"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// CreateOrder creates a new order in the database
func (r *OrderRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	order.ID = uuid.New().String()
	order.OrderedAt = time.Now()
	order.UpdatedAt = time.Now()

	query := `
		INSERT INTO orders (id, user_id, cart_id, order_number, status, payment_status, 
		                    payment_method, shipping_address, billing_address, subtotal,
		                    discount, tax, shipping_fee, total_amount, currency, ordered_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`

	_, err := r.db.Exec(
		ctx,
		query,
		order.ID,
		order.UserID,
		order.CartID,
		order.OrderNumber,
		order.Status,
		order.PaymentStatus,
		order.PaymentMethod,
		order.ShippingAddress,
		order.BillingAddress,
		order.Subtotal,
		order.Discount,
		order.Tax,
		order.ShippingFee,
		order.TotalAmount,
		order.Currency,
		order.OrderedAt,
		order.UpdatedAt)
	return err
}

func (r *OrderRepository) CreateOrderItem(ctx context.Context, item *models.OrderItem) error {
	query := `
        INSERT INTO order_items (
            id, order_id, product_id, product_name, unit_price, quantity, discount, created_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `
	_, err := r.db.Exec(ctx, query,
		uuid.New().String(), item.OrderID, item.ProductID, item.ProductName,
		item.UnitPrice, item.Quantity, item.Discount, item.CreatedAt,
	)
	return err
}

// GetOrderById retrieves an order by its ID
func (r *OrderRepository) GetOrderById(ctx context.Context, id string) (*models.Order, error) {
	query := `
		SELECT id, user_id, cart_id, order_number, status, payment_status, payment_method, shipping_address,
		billing_address, subtotal, discount, tax, shipping_fee, total_amount, currency, ordered_at, updated_at
		FROM orders where id = $1`

	order := &models.Order{}
	row := r.db.QueryRow(ctx, query, id).Scan(
		&order.ID,
		&order.UserID,
		&order.CartID,
		&order.OrderNumber,
		&order.Status,
		&order.PaymentStatus,
		&order.PaymentMethod,
		&order.ShippingAddress,
		&order.BillingAddress,
		&order.Subtotal,
		&order.Discount,
		&order.Tax,
		&order.ShippingFee,
		&order.TotalAmount,
		&order.Currency,
		&order.OrderedAt,
		&order.UpdatedAt)
	if row == pgx.ErrNoRows {
		return nil, nil
	}
	return order, nil
}

// UpdateOrder updates an existing order in the database
func (r *OrderRepository) UpdateOrder(ctx context.Context, order *models.Order) error {
	order.UpdatedAt = time.Now()
	query := `
		UPDATE orders 
		SET status = $2, payment_status = $3, payment_method = $4, shipping_address = $5, billing_address = $6,
		subtotal = $7, discount = $8, tax = $9, shipping_fee = $10, total_amount = $11, currency = $12, ordered_at = $13, 
		updated_at = $14
		WHERE id = $1`
	_, err := r.db.Exec(
		ctx,
		query,
		order.ID,
		order.Status,
		order.PaymentStatus,
		order.PaymentMethod,
		order.ShippingAddress,
		order.BillingAddress,
		order.Subtotal,
		order.Discount,
		order.Tax,
		order.ShippingFee,
		order.TotalAmount,
		order.Currency,
		order.OrderedAt,
		order.UpdatedAt)

	return err
}

// DeleteOrder deletes an order from the database
func (r *OrderRepository) DeleteOrder(ctx context.Context, id string) error {
	query := `
		DELETE FROM orders
		WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// GetOrderItemsByOrderId retrieves all order items by its order number
func (r *OrderRepository) GetOrderItemsByOrderId(ctx context.Context, orderID string) ([]*models.OrderItem, error) {
	query := `
		SELECT id, order_id, product_id, product_name, unit_price, quantity, discount, total_price, created_at
		FROM order_items WHERE order_id = $1`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.OrderItem
	for rows.Next() {
		item := &models.OrderItem{}
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.ProductName,
			&item.UnitPrice,
			&item.Quantity,
			&item.Discount,
			&item.TotalPrice,
			&item.CreatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// GetOrdersByUserID retrieves all orders for a specific user
func (r *OrderRepository) GetOrdersByUserID(ctx context.Context, userID string) ([]models.Order, error) {
	query := `
		SELECT id, user_id, cart_id, order_number, status, payment_status, payment_method, shipping_address,
		billing_address, subtotal, discount, tax, shipping_fee, total_amount, currency, ordered_at, updated_at
		FROM orders WHERE user_id = $1`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		order := models.Order{}
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.CartID,
			&order.OrderNumber,
			&order.Status,
			&order.PaymentStatus,
			&order.PaymentMethod,
			&order.ShippingAddress,
			&order.BillingAddress,
			&order.Subtotal,
			&order.Discount,
			&order.Tax,
			&order.ShippingFee,
			&order.TotalAmount,
			&order.Currency,
			&order.OrderedAt,
			&order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

// GetOrderCountByStatus retrieves the count of orders by their status
func (r *OrderRepository) GetOrderCountByStatus(ctx context.Context, status string) (int, error) {
	query := `
		SELECT COUNT(*) FROM orders WHERE status = $1`

	var count int
	err := r.db.QueryRow(ctx, query, status).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetOrderCountByUserID retrieves the count of orders for a specific user
func (r *OrderRepository) GetOrderCountByUserID(ctx context.Context, userID string) (int, error) {
	query := `
		SELECT COUNT(*) FROM orders WHERE user_id = $1`

	var count int
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Search with filters like status, date range, amount range, etc.
func (r *OrderRepository) FilterOrders(ctx context.Context, filters models.OrderFilter) ([]models.Order, error) {
	query := `
		SELECT id, user_id, cart_id, order_number, status, payment_status, payment_method, shipping_address,
		billing_address, subtotal, discount, tax, shipping_fee, total_amount, currency, ordered_at, updated_at
		FROM orders WHERE 1=1`

	args := []interface{}{}
	argIndex := 1

	if filters.UserID != nil && *filters.UserID != "" {
		query += ` AND user_id = $` + string(argIndex)
		args = append(args, filters.UserID)
		argIndex++
	}
	if filters.Status != nil && *filters.Status != "" {
		query += ` AND status = $` + string(argIndex)
		args = append(args, filters.Status)
		argIndex++
	}
	if filters.PaymentStatus != nil && *filters.PaymentStatus != "" {
		query += ` AND payment_status = $` + string(argIndex)
		args = append(args, filters.PaymentStatus)
		argIndex++
	}
	if filters.StartDate != nil && !filters.StartDate.IsZero() {
		query += ` AND ordered_at >= $` + string(argIndex)
		args = append(args, filters.StartDate)
		argIndex++
	}
	if filters.EndDate != nil && !filters.EndDate.IsZero() {
		query += ` AND ordered_at <= $` + string(argIndex)
		args = append(args, filters.EndDate)
		argIndex++
	}
	if filters.MinTotal != nil {
		query += ` AND total_amount >= $` + string(argIndex)
		args = append(args, filters.MinTotal)
		argIndex++
	}
	if filters.MaxTotal != nil {
		query += ` AND total_amount <= $` + string(argIndex)
		args = append(args, filters.MaxTotal)
		argIndex++
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		order := models.Order{}
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.CartID,
			&order.OrderNumber,
			&order.Status,
			&order.PaymentStatus,
			&order.PaymentMethod,
			&order.ShippingAddress,
			&order.BillingAddress,
			&order.Subtotal,
			&order.Discount,
			&order.Tax,
			&order.ShippingFee,
			&order.TotalAmount,
			&order.Currency,
			&order.OrderedAt,
			&order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

// Support for paginated order listings in admin panels.
func (r *OrderRepository) ListOrders(ctx context.Context, limit int, offset int) ([]models.Order, error) {
	query := `
		SELECT id, user_id, cart_id, order_number, status, payment_status, payment_method, shipping_address,
		billing_address, subtotal, discount, tax, shipping_fee, total_amount, currency, ordered_at, updated_at
		FROM orders
		ORDER BY ordered_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		order := models.Order{}
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.CartID,
			&order.OrderNumber,
			&order.Status,
			&order.PaymentStatus,
			&order.PaymentMethod,
			&order.ShippingAddress,
			&order.BillingAddress,
			&order.Subtotal,
			&order.Discount,
			&order.Tax,
			&order.ShippingFee,
			&order.TotalAmount,
			&order.Currency,
			&order.OrderedAt,
			&order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

// Update only the order status.
func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, orderID string, status string) error {
	query := `
		UPDATE orders 
		SET status = $2, updated_at = $3
		WHERE id = $1`

	_, err := r.db.Exec(ctx, query, orderID, status, time.Now())
	return err
}

// Change payment status (e.g., paid, failed, refunded)
func (r *OrderRepository) UpdatePaymentStatus(ctx context.Context, orderID string, paymentStatus string) error {
	query := `
		UPDATE orders 
		SET payment_status = $2, updated_at = $3
		WHERE id = $1`

	_, err := r.db.Exec(ctx, query, orderID, paymentStatus, time.Now())
	return err
}

// GetOrderCount retrieves the total number of orders in the database
func (r *OrderRepository) GetOrderCount(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*) FROM orders`

	var count int
	err := r.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetOrderCountByDateRange retrieves the count of orders within a specific date range
func (r *OrderRepository) GetOrderCountByDateRange(ctx context.Context, startDate time.Time, endDate time.Time) (int, error) {
	query := `
		SELECT COUNT(*) FROM orders 
		WHERE ordered_at >= $1 AND ordered_at <= $2`

	var count int
	err := r.db.QueryRow(ctx, query, startDate, endDate).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
