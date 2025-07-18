package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/mahabub618/minipack/internal/models"
	"github.com/mahabub618/minipack/internal/repositories"
	"github.com/mahabub618/minipack/internal/services"
	"github.com/mahabub618/minipack/utils"
	"log"
	"net/http"
)

type CartHandler struct {
	cartService       *services.CartService
	anonymousUserRepo *repositories.AnonymousUserRepository
}

func NewCartHandler(cartService *services.CartService, anonymousUserRepo *repositories.AnonymousUserRepository) *CartHandler {
	return &CartHandler{cartService: cartService, anonymousUserRepo: anonymousUserRepo}
}

func (h *CartHandler) CreateCart(w http.ResponseWriter, r *http.Request) {
	var cart models.Cart
	if err := json.NewDecoder(r.Body).Decode(&cart); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	newCart, err := h.cartService.CreateCart(r.Context(), &cart)
	if err != nil {
		http.Error(w, "Failed to create cart: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCart)
}

func (h *CartHandler) GetCartByID(w http.ResponseWriter, r *http.Request) {
	cartId := chi.URLParam(r, "cart_id")
	cart, err := h.cartService.GetCartByID(r.Context(), cartId)
	if err != nil {
		http.Error(w, "Failed to retrieve cart: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if cart == nil {
		http.Error(w, "Cart not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cart)
}

func (h *CartHandler) UpdateCart(w http.ResponseWriter, r *http.Request) {
	cartId := chi.URLParam(r, "cart_id")
	var cart models.Cart
	if err := json.NewDecoder(r.Body).Decode(&cart); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	cart.ID = cartId // Ensure the cart ID is set
	if err := h.cartService.UpdateCart(r.Context(), &cart); err != nil {
		http.Error(w, "Failed to update cart: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cart)
}

func (h *CartHandler) DeleteCart(w http.ResponseWriter, r *http.Request) {
	cartId := chi.URLParam(r, "cart_id")
	if err := h.cartService.DeleteCart(r.Context(), cartId); err != nil {
		http.Error(w, "Failed to delete cart: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CartHandler) GetCartItems(w http.ResponseWriter, r *http.Request) {
	cartId := chi.URLParam(r, "cart_id")
	items, err := h.cartService.GetCartItems(r.Context(), cartId)
	if err != nil {
		http.Error(w, "Failed to retrieve cart items: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}

func (h *CartHandler) AddItemToCart(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetOrCreateAnonymousUserID(w, r, h.anonymousUserRepo)
	// Get or create cart for user
	cart, err := h.cartService.GetActiveCart(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to retrieve cart: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if cart == nil {
		// If no active cart exists, create a new one
		newCart := &models.Cart{
			UserID: userID,
			Status: models.CartStatusActive,
		}
		cart, err = h.cartService.CreateCart(r.Context(), newCart)
		if err != nil {
			http.Error(w, "Failed to create cart: "+err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("No active cart, newly created cart", cart)
	} else {
		log.Println("Found existing active cart", cart)
	}

	// Add item to cart
	var item models.CartItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := h.cartService.AddItemToCart(r.Context(), cart.ID, &item); err != nil {
		http.Error(w, "Failed to add item to cart: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}
