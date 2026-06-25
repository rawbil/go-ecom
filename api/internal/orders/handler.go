package orders

import (
	"encoding/json"
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
			http.Error(w, NotFoundError.Error(), http.StatusNotFound)
		}
		utils.ErrorHandler(err, w, http.StatusInternalServerError)
		return
	}
	utils.JsonResponse(w, createdOrder)

}
