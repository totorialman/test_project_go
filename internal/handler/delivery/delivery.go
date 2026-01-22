package delivery

import (
	"encoding/json"
	"log"
	"net/http"
)

type Handler struct {
	usecase usecase
}

func NewHandler(usecase usecase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var orderID OrderID

	if err := json.NewDecoder(r.Body).Decode(&orderID); err != nil {
		log.Printf(orderID.OrderID)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	delivery, err := h.usecase.Create(r.Context(), orderID.OrderID)
	if err != nil {
		log.Printf("No CREATE: %v", err)
		w.WriteHeader(http.StatusConflict)
		return
	}

	dtoDelivery := DeliveryAssign{
		CourierID: delivery.CourierID,
		OrderID: delivery.OrderID,
		TransportType: delivery.TransportType,
		Deadline: delivery.Deadline,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err:= json.NewEncoder(w).Encode(dtoDelivery); err!=nil{
		log.Printf("JSON encode error: %v", err)
	}
}
