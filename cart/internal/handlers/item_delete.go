package handlers

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"io"
	"net/http"
	"strconv"
)

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
		return
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
}
