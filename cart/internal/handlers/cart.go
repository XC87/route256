package handlers

import (
	"context"
	"net/http"
	"route256.ozon.ru/project/cart/internal/domain"
	"route256.ozon.ru/project/cart/internal/mw"
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
	orderCheckoutURL  = "POST /cart/checkout"
)

type CartService interface {
	AddItem(userId int64, item domain.Item) error
	GetItemsByUserId(userId int64) (*service.CartResponse, error)
	DeleteItem(userId int64, skuId int64) error
	DeleteItemsByUserId(userId int64) error
	OrderCheckout(ctx context.Context, userId int64) (int64, error)
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
	chain := []mw.MiddlewareChain{mw.LoggingMiddleware}

	http.Handle(addToCartURL, mw.BuildMiddleware(h.AddItem, chain))
	http.Handle(deleteFromCartURL, mw.BuildMiddleware(h.DeleteItem, chain))
	http.Handle(deleteCartURL, mw.BuildMiddleware(h.DeleteItemsByUserId, chain))
	http.Handle(getCartURL, mw.BuildMiddleware(h.GetItemsByUserId, chain))
	http.Handle(orderCheckoutURL, mw.BuildMiddleware(h.OrderCheckout, chain))
}
