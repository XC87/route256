package handlers

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func (h *Handler) GetItemsByUserId(w http.ResponseWriter, r *http.Request) {
	userId, _ := strconv.ParseInt(r.PathValue(userIdPath), 10, 64)
	cartGetReq := CartGetRequest{UserId: userId}
	if _, err := govalidator.ValidateStruct(cartGetReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userCart, err := h.cartService.GetItemsByUserId(r.Context(), cartGetReq.UserId)
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
		zap.L().Info("something went wrong on writing the response:", zap.Error(err))
		return
	}
}
