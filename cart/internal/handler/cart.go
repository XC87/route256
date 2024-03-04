package handler

import (
	"encoding/json"
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/justinas/alice"
	"io"
	"log"
	"net/http"
	"route256.ozon.ru/project/cart/internal/client/product"
	"route256.ozon.ru/project/cart/internal/middleware"
	"strconv"
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
	AddItem(userId int64, skuId int64, count uint16) error
	GetItemsByUserId(userId int64) (*CartResponse, error)
	DeleteItem(userId int64, skuId int64) error
	DeleteItemsByUserId(userId int64) error
}

func NewCartHandler(cartService CartService) *Handler {
	return &Handler{
		cartService: cartService,
	}
}

type Handler struct {
	cartService CartService
}

func (h *Handler) Register() {
	chain := alice.New(middleware.LoggingMiddleware)

	http.Handle(addToCartURL, chain.ThenFunc(h.AddToCart))
	http.Handle(deleteFromCartURL, chain.ThenFunc(h.DeleteItem))
	http.Handle(deleteCartURL, chain.ThenFunc(h.DeleteItemsByUserId))
	http.Handle(getCartURL, chain.ThenFunc(h.GetItemsByUserId))
}

func (h *Handler) AddToCart(w http.ResponseWriter, r *http.Request) {
	userId, _ := strconv.ParseInt(r.PathValue(userIdPath), 10, 64)
	skuId, _ := strconv.ParseInt(r.PathValue(skuIdPath), 10, 64)

	itemAdd := ItemAddRequest{UserId: userId, SkuId: skuId}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "please try again later", http.StatusInternalServerError)
	}

	if err = json.Unmarshal(data, &itemAdd); err != nil {
		http.Error(w, "bad request arguments", http.StatusBadRequest)
	}

	if _, err = govalidator.ValidateStruct(itemAdd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.cartService.AddItem(itemAdd.UserId, itemAdd.SkuId, itemAdd.Count)
	if errors.Is(err, product.ErrProductNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "can't add to cart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	userId, _ := strconv.ParseInt(r.PathValue(userIdPath), 10, 64)
	skuId, _ := strconv.ParseInt(r.PathValue(skuIdPath), 10, 64)

	itemDeleteReq := ItemDeleteRequest{UserId: userId, SkuId: skuId}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "please try again later", http.StatusInternalServerError)
	}

	if err = json.Unmarshal(data, &itemDeleteReq); err != nil {
		http.Error(w, "Bad request arguments", http.StatusBadRequest)
	}

	if _, err = govalidator.ValidateStruct(itemDeleteReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.cartService.DeleteItem(itemDeleteReq.UserId, itemDeleteReq.SkuId)
	if err != nil {
		http.Error(w, "cant delete from cart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	return
}

func (h *Handler) DeleteItemsByUserId(w http.ResponseWriter, r *http.Request) {
	userId, _ := strconv.ParseInt(r.PathValue(userIdPath), 10, 64)
	cartClearReq := CartClearRequest{UserId: userId}
	if _, err := govalidator.ValidateStruct(cartClearReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.cartService.DeleteItemsByUserId(cartClearReq.UserId); err != nil {
		http.Error(w, "cant delete cart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	return
}

func (h *Handler) GetItemsByUserId(w http.ResponseWriter, r *http.Request) {
	userId, _ := strconv.ParseInt(r.PathValue(userIdPath), 10, 64)
	cartGetReq := CartGetRequest{UserId: userId}
	if _, err := govalidator.ValidateStruct(cartGetReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userCart, err := h.cartService.GetItemsByUserId(cartGetReq.UserId)
	jsonResponse := []byte("{}")
	status := http.StatusNotFound

	if err == nil && len(userCart.Items) > 0 {
		jsonResponse, _ = json.Marshal(userCart)
		status = http.StatusOK
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Fatalln("something when wrong on write answer")
		return
	}
	return
}
