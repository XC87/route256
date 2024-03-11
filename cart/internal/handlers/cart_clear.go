package handlers

import (
	"github.com/asaskevich/govalidator"
	"net/http"
	"strconv"
)

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
}
