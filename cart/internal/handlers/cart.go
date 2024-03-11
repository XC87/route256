package handlers

import (
	"net/http"
	"route256.ozon.ru/project/cart/internal/domain"
	"route256.ozon.ru/project/cart/internal/service"
)

const (
	userIdPath = "user_id"
	skuIdPath  = "sku_id"
)
const (
	addToCartURL      = "POST /user/{user_id}/cart/{sku_id}"
	deleteFromCartURL = "DELETE /user/{user_id}/cart/{sku_id}"
	deleteCartURL     = "DELETE /user/{user_id}/cart"
	getCartURL        = "GET /user/{user_id}/cart"
)

type CartService interface {
	AddItem(userId int64, item domain.Item) error
	GetItemsByUserId(userId int64) (*service.CartResponse, error)
	DeleteItem(userId int64, skuId int64) error
	DeleteItemsByUserId(userId int64) error
}

type Handler struct {
	cartService CartService
}

func NewCartHandler(cartService CartService) *Handler {
	return &Handler{
		cartService: cartService,
	}
}

func (h *Handler) Register() {
	chain := []middlewareChain{loggingMiddleware}

	http.Handle(addToCartURL, buildMiddleware(h.AddItem, chain))
	http.Handle(deleteFromCartURL, buildMiddleware(h.DeleteItem, chain))
	http.Handle(deleteCartURL, buildMiddleware(h.DeleteItemsByUserId, chain))
	http.Handle(getCartURL, buildMiddleware(h.GetItemsByUserId, chain))
}
