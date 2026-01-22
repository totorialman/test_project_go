package courier

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	courierErrors "github.com/totorialman/test_project_go/internal/errors/courier"
	"github.com/totorialman/test_project_go/internal/usecase/courier"
)

type Handler struct {
	usecase usecase
}

func NewHandler(usecase usecase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var courierDTO CreateCourierDTO

	if err := json.NewDecoder(r.Body).Decode(&courierDTO); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(courierDTO)
	if len(courierDTO.Name) == 0 || len(courierDTO.Name) > 100 || len(courierDTO.Phone) == 0 || len(courierDTO.Status) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	courier := courier.Courier{
		Name:          courierDTO.Name,
		Phone:         courierDTO.Phone,
		Status:        courierDTO.Status,
		TransportType: courierDTO.TransportType,
	}

	err := h.usecase.Create(r.Context(), courier)
	if errors.Is(err, courierErrors.ErrCourierExists) {
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	couriers, err := h.usecase.GetAll(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dtoCouriers := make([]GetCouriersDTO, len(couriers))
	for i, c := range couriers {
		dtoCouriers[i] = GetCouriersDTO{
			ID:            c.ID,
			Name:          c.Name,
			Phone:         c.Phone,
			Status:        c.Status,
			TransportType: c.TransportType,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dtoCouriers); err != nil {
		log.Printf("JSON encode error: %v", err)
	}
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	courier, err := h.usecase.GetById(r.Context(), id)
	if errors.Is(err, courierErrors.ErrCourierNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dtoCourier := GetCouriersDTO{
		ID:            courier.ID,
		Name:          courier.Name,
		Phone:         courier.Phone,
		Status:        courier.Status,
		TransportType: courier.TransportType,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dtoCourier); err != nil {
		log.Printf("JSON encode error: %v", err)
	}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var dtoCourier CreateCourierDTO

	if err := json.NewDecoder(r.Body).Decode(&dtoCourier); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if dtoCourier.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	courier := courier.Courier{
		ID:            dtoCourier.ID,
		Name:          dtoCourier.Name,
		Phone:         dtoCourier.Phone,
		Status:        dtoCourier.Status,
		TransportType: dtoCourier.TransportType,
	}

	err := h.usecase.Update(r.Context(), courier)
	if errors.Is(err, courierErrors.ErrCourierNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

}
