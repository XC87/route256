package handlers

import (
	"encoding/json"
	"errors"
	"github.com/asaskevich/govalidator"
	"net/http"
	"route256.ozon.ru/project/cart/internal/domain"
	"route256.ozon.ru/project/cart/internal/service"
	"strconv"
)

func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	userId, _ := strconv.ParseInt(r.PathValue(userIdPath), 10, 64)
	skuId, _ := strconv.ParseInt(r.PathValue(skuIdPath), 10, 64)

	itemAdd := ItemAddRequest{UserId: userId, SkuId: skuId}
	if err := json.NewDecoder(r.Body).Decode(&itemAdd); err != nil {
		http.Error(w, "bad request arguments", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if _, err := govalidator.ValidateStruct(itemAdd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := h.cartService.AddItem(r.Context(), itemAdd.UserId, domain.Item{
		Sku_id: itemAdd.SkuId,
		Count:  itemAdd.Count,
	}); err != nil {
		switch {
		case errors.Is(err, service.ErrProductNotFound):
			w.WriteHeader(http.StatusBadRequest)
		case errors.Is(err, service.ErrStockInsufficient):
			w.WriteHeader(http.StatusPreconditionFailed)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
