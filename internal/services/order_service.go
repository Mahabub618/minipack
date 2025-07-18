package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/mahabub618/minipack/internal/models"
	"github.com/mahabub618/minipack/internal/repositories"
	"log"
	"time"
)

type OrderService struct {
	orderRepo *repositories.OrderRepository
	cartRepo  *repositories.CartRepository
}

func NewOrderService(orderRepo *repositories.OrderRepository, cartRepo *repositories.CartRepository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		cartRepo:  cartRepo,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *models.Order) error {
	if order.CartID == nil {
		cart, err := s.cartRepo.GetCartByID(ctx, *order.CartID)
		if err != nil {
			return err
		}
		if cart == nil {
			return errors.New("cart not found")
		}
	}
	return s.orderRepo.CreateOrder(ctx, order)
}

func (s *OrderService) GetOrderByID(ctx context.Context, id string) (*models.Order, error) {
	order, err := s.orderRepo.GetOrderById(ctx, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	return order, nil
}

func (s *OrderService) UpdateOrder(ctx context.Context, id string, order *models.Order) error {
	orderItem, err := s.orderRepo.GetOrderById(ctx, id)
	if err != nil {
		return err
	}
	if orderItem == nil {
		return errors.New("order not found")
	}
	return s.orderRepo.UpdateOrder(ctx, order)
}

func (s *OrderService) CreateOrderFromCart(ctx context.Context, cartID string, shippingAddress, billingAddress map[string]interface{}) (*models.Order, error) {
	// Start a transaction
	tx, err := s.orderRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Get cart with items
	cart, err := s.cartRepo.GetCartByID(ctx, cartID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, errors.New("cart not found")
	}
	if cart.Status != models.CartStatusActive {
		return nil, errors.New("cart is not active")
	}

	cartItems, err := s.cartRepo.GetCartItems(ctx, cartID)
	if err != nil {
		return nil, err
	}
	if len(cartItems) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Create order from cart
	order := NewOrderFromCart(cart, cartItems, shippingAddress, billingAddress)
	if err := s.orderRepo.CreateOrder(ctx, order); err != nil {
		return nil, err
	}

	log.Println("Created order: ", order.ID, " || Order Number: ", order.OrderNumber)

	// Convert cart items to order items
	for _, cartItem := range cartItems {
		orderItem := &models.OrderItem{
			OrderID:     order.ID,
			ProductID:   cartItem.ProductID,
			ProductName: cartItem.ProductName,
			UnitPrice:   cartItem.UnitPrice,
			Quantity:    cartItem.Quantity,
			Discount:    cartItem.Discount,
			TotalPrice:  cartItem.TotalPrice,
			CreatedAt:   cartItem.CreatedAt,
		}
		log.Println("OrderId: ", orderItem.OrderID, " || ProductID: ", orderItem.ProductID, " || Quantity: ", orderItem.Quantity,
			" || TotalPrice: ", orderItem.TotalPrice, " || Discount: ", orderItem.Discount, " || UnitPrice: ", orderItem.UnitPrice)
		if err := s.orderRepo.CreateOrderItem(ctx, orderItem); err != nil {
			return nil, err
		}
	}

	// Update cart status to converted
	cart.Status = models.CartStatusConverted
	if err := s.cartRepo.UpdateCart(ctx, cart); err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return order, nil
}

// NewOrderFromCart creates a new order from a cart
func NewOrderFromCart(cart *models.Cart, cartItems []models.CartItem, shippingAddress, billingAddress map[string]interface{}) *models.Order {
	now := time.Now()

	// Calculate totals
	subtotal := 0.0
	for _, item := range cartItems {
		subtotal += item.TotalPrice
	}

	tax := subtotal * 0.10 // Example: 10% tax
	shippingFee := 5.99    // Example flat shipping fee

	return &models.Order{
		UserID:          cart.UserID,
		CartID:          &cart.ID,
		OrderNumber:     generateOrderNumber(),
		Status:          models.OrderStatusPending,
		PaymentStatus:   models.PaymentStatusPending,
		ShippingAddress: shippingAddress,
		BillingAddress:  billingAddress,
		Subtotal:        subtotal,
		Tax:             tax,
		ShippingFee:     shippingFee,
		TotalAmount:     subtotal + tax + shippingFee,
		Currency:        "USD",
		OrderedAt:       now,
		UpdatedAt:       now,
	}
}

func generateOrderNumber() string {
	orderNumber := "ORD-" + time.Now().Format("20060102") + "-" + uuid.New().String()[:8]
	log.Println("#255 generateOrderNumber: ", orderNumber, " || Length: ", len(orderNumber))
	return orderNumber
}

func (s *OrderService) DeleteOrder(ctx context.Context, id string) error {
	order, err := s.orderRepo.GetOrderById(ctx, id)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("order not found")
	}

	return s.orderRepo.DeleteOrder(ctx, id)
}

func (s *OrderService) GetOrderItems(ctx context.Context, orderID string) ([]*models.OrderItem, error) {
	items, err := s.orderRepo.GetOrderItemsByOrderId(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if items == nil {
		return nil, errors.New("no items found for this order")
	}
	return items, nil
}

func (s *OrderService) GetOrdersByUserID(ctx context.Context, userID string) ([]models.Order, error) {
	orders, err := s.orderRepo.GetOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if orders == nil {
		return nil, errors.New("no orders found for this user")
	}
	return orders, nil
}

func (s *OrderService) GetOrderCountByStatus(ctx context.Context, status string) (int, error) {
	count, err := s.orderRepo.GetOrderCountByStatus(ctx, status)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *OrderService) GetOrderCountByUserID(ctx context.Context, userID string) (int, error) {
	count, err := s.orderRepo.GetOrderCountByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *OrderService) FilterOrders(ctx context.Context, filters models.OrderFilter) ([]models.Order, error) {
	orders, err := s.orderRepo.FilterOrders(ctx, filters)
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, errors.New("no orders found matching the filters")
	}
	return orders, nil
}

func (s *OrderService) ListOrders(ctx context.Context, limit, offset int) ([]models.Order, error) {
	orders, err := s.orderRepo.ListOrders(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, errors.New("no orders found")
	}
	return orders, nil
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID, status string) error {
	order, err := s.orderRepo.GetOrderById(ctx, orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("order not found")
	}
	return s.orderRepo.UpdateOrderStatus(ctx, orderID, status)
}

func (s *OrderService) UpdatePaymentStatus(ctx context.Context, orderID, paymentStatus string) error {
	order, err := s.orderRepo.GetOrderById(ctx, orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("order not found")
	}
	return s.orderRepo.UpdatePaymentStatus(ctx, orderID, paymentStatus)
}

func (s *OrderService) GetOrderCount(ctx context.Context) (int, error) {
	count, err := s.orderRepo.GetOrderCount(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *OrderService) GetOrderCountByDateRange(ctx context.Context, startDate, endDate string) (int, error) {
	// how to convert startDate and endDate to time.Time?
	if startDate == "" || endDate == "" {
		return 0, errors.New("start date and end date cannot be empty")
	}
	// Assuming startDate and endDate are in the format "YYYY-MM-DD"
	startDateTime, err := time.Parse("YYYY-MM-DD", startDate)
	if err != nil {
		return 0, errors.New("invalid start date format")
	}
	endDateTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return 0, errors.New("invalid end date format")
	}
	if startDateTime.After(endDateTime) {
		return 0, errors.New("start date cannot be after end date")
	}

	count, err := s.orderRepo.GetOrderCountByDateRange(ctx, startDateTime, endDateTime)
	if err != nil {
		return 0, err
	}
	return count, nil
}
