package handlers

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"log"
	"net/http"
)

func (h *Handler) OrderCheckout(w http.ResponseWriter, r *http.Request) {
	cartGetReq := OrderCheckoutRequest{}
	if err := json.NewDecoder(r.Body).Decode(&cartGetReq); err != nil {
		http.Error(w, "bad request arguments", http.StatusBadRequest)
		return
	}
	if _, err := govalidator.ValidateStruct(cartGetReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	orderId, err := h.cartService.OrderCheckout(r.Context(), cartGetReq.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusPreconditionFailed)
		return
	}

	jsonResponse, _ := json.Marshal(&OrderCheckoutResponse{orderId})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Println("error writing response:", err)
		return
	}
}
