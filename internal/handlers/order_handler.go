package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/mahabub618/minipack/internal/models"
	"github.com/mahabub618/minipack/internal/repositories"
	"github.com/mahabub618/minipack/internal/services"
	"github.com/mahabub618/minipack/utils"
	"net/http"
	"strconv"
	"time"
)

type OrderHandler struct {
	orderService      *services.OrderService
	cartService       *services.CartService
	anonymousUserRepo *repositories.AnonymousUserRepository
}

func NewOrderHandler(orderService *services.OrderService, cartService *services.CartService, anonymousUserRepo *repositories.AnonymousUserRepository) *OrderHandler {
	return &OrderHandler{orderService: orderService, cartService: cartService, anonymousUserRepo: anonymousUserRepo}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := h.orderService.CreateOrder(r.Context(), &order); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	order, err := h.orderService.GetOrderByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if order == nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var order map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get the existing order to ensure it exists
	existingOrder, err := h.orderService.GetOrderByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to get platform: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if existingOrder == nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	for key, value := range order {
		switch key {
		case "status":
			existingOrder.Status = models.OrderStatus(value.(string))
		case "payment_status":
			existingOrder.PaymentStatus = models.PaymentStatus(value.(string))
		case "payment_method":
			if paymentMethod, ok := value.(string); ok {
				existingOrder.PaymentMethod = &paymentMethod
			}
		case "shipping_address":
			if shippingAddress, ok := value.(map[string]interface{}); ok {
				existingOrder.ShippingAddress = shippingAddress
			}
		case "billing_address":
			if billingAddress, ok := value.(map[string]interface{}); ok {
				existingOrder.BillingAddress = billingAddress
			}
		case "subtotal":
			if subtotal, ok := value.(float64); ok {
				existingOrder.Subtotal = subtotal
			}
		case "discount":
			if discount, ok := value.(float64); ok {
				existingOrder.Discount = discount
			}
		case "tax":
			if tax, ok := value.(float64); ok {
				existingOrder.Tax = tax
			}
		case "shipping_fee":
			if shippingFee, ok := value.(float64); ok {
				existingOrder.ShippingFee = shippingFee
			}
		case "total_amount":
			if totalAmount, ok := value.(float64); ok {
				existingOrder.TotalAmount = totalAmount
			}
		case "currency":
			if currency, ok := value.(string); ok {
				existingOrder.Currency = currency
			}
		case "ordered_at":
			if orderedAt, ok := value.(string); ok {
				parsedTime, err := time.Parse(time.RFC3339, orderedAt)
				if err != nil {
					http.Error(w, "Invalid ordered_at format", http.StatusBadRequest)
					return
				}
				existingOrder.OrderedAt = parsedTime
			}
		}
	}

	if err := h.orderService.UpdateOrder(r.Context(), id, existingOrder); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingOrder)
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.orderService.DeleteOrder(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *OrderHandler) GetOrderItems(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")
	items, err := h.orderService.GetOrderItems(r.Context(), orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if items == nil {
		http.Error(w, "No items found for this order", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}

func (h *OrderHandler) GetOrdersByUserID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	orders, err := h.orderService.GetOrdersByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if orders == nil {
		http.Error(w, "No orders found for this user", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) FilterOrders(w http.ResponseWriter, r *http.Request) {
	var filters models.OrderFilter
	if err := json.NewDecoder(r.Body).Decode(&filters); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	orders, err := h.orderService.FilterOrders(r.Context(), filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if orders == nil {
		http.Error(w, "No orders found for the given filters", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))
	orders, err := h.orderService.ListOrders(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if orders == nil {
		http.Error(w, "No orders found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var payload struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := h.orderService.UpdateOrderStatus(r.Context(), id, payload.Status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *OrderHandler) UpdatePaymentStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var payload struct {
		PaymentStatus string `json:"payment_status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := h.orderService.UpdatePaymentStatus(r.Context(), id, payload.PaymentStatus); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"payment_status": "updated"})
}

func (h *OrderHandler) GetOrderCount(w http.ResponseWriter, r *http.Request) {
	count, err := h.orderService.GetOrderCount(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"count": count})
}

func (h *OrderHandler) GetOrderCountByDateRange(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	startDate := query.Get("start_date")
	endDate := query.Get("end_date")
	count, err := h.orderService.GetOrderCountByDateRange(r.Context(), startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"count": count})
}

// CreateOrderFromCart creates an order from the user's active cart
func (h *OrderHandler) CreateOrderFromCart(w http.ResponseWriter, r *http.Request) {
	// Get user ID (authenticated or from session)
	userID, err := utils.GetOrCreateAnonymousUserID(w, r, h.anonymousUserRepo)
	// Get active cart
	cart, err := h.cartService.GetActiveCart(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to retrieve cart: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if cart == nil {
		http.Error(w, "No active cart found for user", http.StatusNotFound)
		return
	}

	// Parse shipping/billing info
	var request struct {
		ShippingAddress map[string]interface{} `json:"shipping_address"`
		BillingAddress  map[string]interface{} `json:"billing_address"`
		PaymentMethod   string                 `json:"payment_method"`
	}
	// Parse request...
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Create order
	order, err := h.orderService.CreateOrderFromCart(r.Context(),
		cart.ID,
		request.ShippingAddress,
		request.BillingAddress)
	if err != nil {
		http.Error(w, "Failed to create order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Process payment (simplified)
	//paymentResult, err := h.paymentService.ProcessPayment(r.Context(),
	//	order.TotalAmount,
	//	request.PaymentMethod)
	//if err != nil {
	//	// Handle payment error
	//	h.orderService.UpdateOrderStatus(r.Context(), order.ID, models.OrderStatusPaymentFailed)
	//	// Return error
	//}

	// Update order status
	//if paymentResult.Success {
	//	h.orderService.UpdateOrderStatus(r.Context(), order.ID, models.OrderStatusPaid)
	//	h.orderService.UpdatePaymentStatus(r.Context(), order.ID, models.PaymentStatusPaid)
	//}

	// Return order confirmation
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}
