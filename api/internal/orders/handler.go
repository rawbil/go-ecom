package orders

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/rawbil/ecom2/internal/utils"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// ! CREATEORDER
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var params CreateOrderParams

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}

	createdOrder, err := h.service.CreateOrder(r.Context(), params)

	if err != nil {

		if err == NotFoundError {
			utils.ErrorHandler(err, err.Error(), w, http.StatusNotFound)
			return
		}
		utils.ErrorHandler(err, "oops... server error", w, http.StatusInternalServerError)
		return
	}
	utils.JsonResponse(w, createdOrder)

}

// ! GETMYORDER
func (h *Handler) GetMyOrder(w http.ResponseWriter, r *http.Request) {

	order_details, err := h.service.GetMyOrder(r.Context())
	if err != nil {
		if err == AuthNotFound {
			utils.ErrorHandler(err, err.Error(), w, http.StatusUnauthorized)
			return
		}
		utils.ErrorHandler(err, "oops... server error", w, http.StatusInternalServerError)
		return
	}

	if len(order_details) < 1 {
		utils.ErrorHandler(errors.New("No orders for you"), "No orders for you", w, http.StatusNotFound)
		return
	}

	utils.JsonResponse(w, utils.SuccessMessage{
		Message: "orders fetched successfully",
		Data:    map[string]any{"orders": order_details},
	})

}

// ! GetAllOrders
func (h *Handler) GetAllOrders(w http.ResponseWriter, r *http.Request) {

	order_details, err := h.service.GetAllOrders(r.Context())
	if err != nil {
		utils.ErrorHandler(err, "oops... server error", w, http.StatusInternalServerError)
		return
	}

	if len(order_details) < 1 {
		utils.ErrorHandler(errors.New("No orders found"), "No orders found", w, http.StatusNotFound)
		return
	}

	utils.JsonResponse(w, utils.SuccessMessage{
		Message: "orders fetched successfully",
		Data:    map[string]any{"orders": order_details},
	})
}

// ! CancelOrder
func (h *Handler) CancleOrder(w http.ResponseWriter, r *http.Request) {
	type Body struct {
		OrderID int64 `json:"order_id"`
	}
	var body Body
	if err := utils.DecodeClient(r, &body); err != nil {
		utils.ErrorHandler(err, "error decoding body", w, http.StatusBadRequest)
		return
	}

	if _, err := h.service.CancelOrder(r.Context(), body.OrderID); err != nil {
		if err == AuthNotFound {
			utils.ErrorHandler(err, err.Error(), w, http.StatusUnauthorized)
			return
		}
		utils.ErrorHandler(err, "oops... server error", w, http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, utils.SuccessMessage{
		Message: "order cancelled",
	})
}
