package handlers

import (
	"context"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
	"route256.ozon.ru/pkg/metrics"
	"route256.ozon.ru/pkg/tracer"
	"route256.ozon.ru/project/cart/internal/domain"
	"route256.ozon.ru/project/cart/internal/mw"
	"route256.ozon.ru/project/cart/internal/service"
)

const (
	userIdPath = "user_id"
	skuIdPath  = "sku_id"
)
const (
	addItem             = "POST /user/{user_id}/cart/{sku_id}"
	deleteItem          = "DELETE /user/{user_id}/cart/{sku_id}"
	deleteItemsByUserId = "DELETE /user/{user_id}/cart"
	getItemsByUserId    = "GET /user/{user_id}/cart"
	orderCheckout       = "POST /cart/checkout"
)

type CartService interface {
	AddItem(ctx context.Context, userId int64, item domain.Item) error
	GetItemsByUserId(ctx context.Context, userId int64) (*service.CartResponse, error)
	DeleteItem(ctx context.Context, userId int64, skuId int64) error
	DeleteItemsByUserId(ctx context.Context, userId int64) error
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

func (h *Handler) Register(serviceName string) {
	mainChain := []mw.MiddlewareChain{tracer.HandleMiddleware, mw.LoggingMiddleware}
	// возможно я тут сильно загнался, но цель была в том чтобы было полноценное описание
	registerURL(serviceName, h.AddItem, "addItem", addItem, mainChain)
	registerURL(serviceName, h.DeleteItem, "deleteItem", deleteItem, mainChain)
	registerURL(serviceName, h.DeleteItemsByUserId, "deleteItemsByUserId", deleteItemsByUserId, mainChain)
	registerURL(serviceName, h.GetItemsByUserId, "getItemsByUserId", getItemsByUserId, mainChain)
	registerURL(serviceName, h.OrderCheckout, "orderCheckout", orderCheckout, mainChain)
}

func registerURL(serviceName string, handler http.HandlerFunc, endpointName string, endpointURL string, mainChain []mw.MiddlewareChain) {
	mwChain := append(mainChain, metrics.CreateMetricMiddleware(serviceName, endpointName, endpointURL))
	handler = mw.BuildMiddleware(handler, mwChain)
	http.Handle(endpointURL, otelhttp.NewHandler(handler, endpointName))
}
