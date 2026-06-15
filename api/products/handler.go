package products

import (
	"encoding/json"
	"net/http"
)

type handler struct {
	service Service
}

type Products struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	// products := []Products{
	// 	{1, "Pampers", 20},
	// 	{2, "Bulb", 30},
	// }

	products, err := h.service.GetAllProducts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(products)
}
