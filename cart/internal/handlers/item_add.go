package handlers

import (
	"encoding/json"
	"errors"
	"github.com/asaskevich/govalidator"
	"io"
	"net/http"
	"route256.ozon.ru/project/cart/internal/domain"
	"route256.ozon.ru/project/cart/internal/service"
	"strconv"
)

func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	userId, _ := strconv.ParseInt(r.PathValue(userIdPath), 10, 64)
	skuId, _ := strconv.ParseInt(r.PathValue(skuIdPath), 10, 64)

	itemAdd := ItemAddRequest{UserId: userId, SkuId: skuId}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "please try again later", http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(data, &itemAdd); err != nil {
		http.Error(w, "bad request arguments", http.StatusBadRequest)
		return
	}

	if _, err = govalidator.ValidateStruct(itemAdd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.cartService.AddItem(itemAdd.UserId, domain.Item{
		Sku_id: itemAdd.SkuId,
		Count:  itemAdd.Count,
	})
	if errors.Is(err, service.ErrProductNotFound) {
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
}
