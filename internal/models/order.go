package models

import (
	"time"
)

type OrderStatus string
type PaymentStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusCompleted OrderStatus = "completed"

	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusPaid    PaymentStatus = "paid"
	PaymentStatusFailed  PaymentStatus = "failed"
	PaymentStatusRefund  PaymentStatus = "refunded"
)

type Order struct {
	ID            string        `json:"id"`
	UserID        string        `json:"user_id"`
	CartID        *string       `json:"cart_id,omitempty"`
	OrderNumber   string        `json:"order_number"`
	Status        OrderStatus   `json:"status"`
	PaymentStatus PaymentStatus `json:"payment_status"`
	PaymentMethod *string       `json:"payment_method,omitempty"`

	ShippingAddress map[string]interface{} `json:"shipping_address"`
	BillingAddress  map[string]interface{} `json:"billing_address"`

	Subtotal    float64 `json:"subtotal"`
	Discount    float64 `json:"discount"`
	Tax         float64 `json:"tax"`
	ShippingFee float64 `json:"shipping_fee"`
	TotalAmount float64 `json:"total_amount"`
	Currency    string  `json:"currency"`

	OrderedAt time.Time `json:"ordered_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrderItem struct {
	ID          string    `json:"id"`
	OrderID     string    `json:"order_id"`
	ProductID   string    `json:"product_id"`
	ProductName string    `json:"product_name"`
	UnitPrice   float64   `json:"unit_price"`
	Quantity    int       `json:"quantity"`
	Discount    float64   `json:"discount"`
	TotalPrice  float64   `json:"total_price"`
	CreatedAt   time.Time `json:"created_at"`
}
