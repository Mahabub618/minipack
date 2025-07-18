package services

import (
	"context"
	"errors"
	"github.com/mahabub618/minipack/internal/models"
	"github.com/mahabub618/minipack/internal/repositories"
)

type CartService struct {
	cartRepo  *repositories.CartRepository
	orderRepo *repositories.OrderRepository
}

func NewCartService(cartRepo *repositories.CartRepository, orderRepo *repositories.OrderRepository) *CartService {
	return &CartService{
		cartRepo:  cartRepo,
		orderRepo: orderRepo,
	}
}

func (s *CartService) CreateCart(ctx context.Context, cart *models.Cart) (*models.Cart, error) {
	err := s.cartRepo.CreateCart(ctx, cart)
	if err != nil {
		return nil, err
	}
	return cart, nil
}

func (s *CartService) GetCartByID(ctx context.Context, id string) (*models.Cart, error) {
	cart, err := s.cartRepo.GetCartByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, nil // or return an error if preferred
	}
	return cart, nil
}

func (s *CartService) UpdateCart(ctx context.Context, cart *models.Cart) error {
	return s.cartRepo.UpdateCart(ctx, cart)
}

func (s *CartService) DeleteCart(ctx context.Context, id string) error {
	return s.cartRepo.DeleteCart(ctx, id)
}

func (s *CartService) GetCartItems(ctx context.Context, cartID string) ([]models.CartItem, error) {
	items, err := s.cartRepo.GetCartItems(ctx, cartID)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *CartService) AddItemToCart(ctx context.Context, cartID string, item *models.CartItem) error {
	cart, err := s.cartRepo.GetCartByID(ctx, cartID)
	if err != nil {
		return err
	}
	if cart == nil {
		return errors.New("cart not found")
	}
	if cart.Status != models.CartStatusActive {
		return errors.New("cannot add items to a non-active cart")
	}

	item.CartID = cartID
	return s.cartRepo.AddItemToCart(ctx, item)
}

func (s *CartService) RemoveItemFromCart(ctx context.Context, cartID, itemID string) error {
	cart, err := s.cartRepo.GetCartByID(ctx, cartID)
	if err != nil {
		return err
	}
	if cart == nil {
		return errors.New("cart not found")
	}
	if cart.Status != models.CartStatusActive {
		return errors.New("cannot remove items from a non-active cart")
	}
	return s.cartRepo.RemoveItemFromCart(ctx, itemID)
}

func (s *CartService) UpdateCartItemQuantity(ctx context.Context, cartID, itemID string, quantity int) error {
	cart, err := s.cartRepo.GetCartByID(ctx, cartID)
	if err != nil {
		return err
	}
	if cart == nil {
		return errors.New("cart not found")
	}
	if cart.Status != models.CartStatusActive {
		return errors.New("cannot update items in a non-active cart")
	}
	return s.cartRepo.UpdateCartItemQuantity(ctx, itemID, quantity)
}

func (s *CartService) GetActiveCart(ctx context.Context, userID string) (*models.Cart, error) {
	cart, err := s.cartRepo.GetActiveCart(ctx, userID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, nil // or return an error if preferred
	}
	return cart, nil
}
