package products

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	repository "github.com/rawbil/ecom2/internal/adapters/sqlc"
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

// ! CreateProduct
func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var params repository.CreateProductParams

	json.NewDecoder(r.Body).Decode(&params)

	result, err := h.service.CreateProduct(r.Context(), params)
	if err != nil {
		utils.ErrorHandler(err, w, http.StatusInternalServerError)
		return
	}

	product_id, _ := result.LastInsertId()

	product, _ := h.service.ListProduct(r.Context(), product_id)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s created successfully!!", product.ProductName)

}

// ! ListProducts
func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	// Get query from request
	query := r.URL.Query()

	page := query.Get("page")
	limit := query.Get("limit")

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 || limitInt > 30 {
		pageInt = 1
	}

	product_name := query.Get("name")
	min_price := query.Get("min_price")
	max_price := query.Get("max_price")

	minPrice, err := strconv.Atoi(min_price)
	if err != nil || minPrice < 1 {
		minPrice = 0
	}
	maxPrice, err := strconv.Atoi(max_price)
	if err != nil || maxPrice < 1 {
		maxPrice = 0
	}

	offset := limitInt * (pageInt - 1)

	products, err := h.service.ListProducts(r.Context(), repository.ListProductsParams{
		Name:     product_name,
		MinPrice: int64(minPrice),
		MaxPrice: int64(maxPrice),
		Limit:    int32(limitInt),
		Offset:   int32(offset),
	})
	if err != nil {
		utils.ErrorHandler(err, w, http.StatusInternalServerError)
		return
	}

	if len(products) < 1 {
		http.Error(w, "No product found", http.StatusNotFound)
		fmt.Println("No products found")
		return
	}

	utils.JsonResponse(w, products)
}

// ! ListProduct
func (h *Handler) ListProduct(w http.ResponseWriter, r *http.Request) {
	type Body struct {
		ID int64
	}
	var body Body
	json.NewDecoder(r.Body).Decode(&body)
	product, err := h.service.ListProduct(r.Context(), body.ID)
	if err != nil {
		utils.ErrorHandler(err, w, http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, product)
}

// ! DeleteProduct
func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	type Body struct {
		ID int64
	}
	var body Body
	json.NewDecoder(r.Body).Decode(&body)
	err := h.service.DeleteProduct(r.Context(), body.ID)
	if err != nil {
		utils.ErrorHandler(err, w, http.StatusInternalServerError)
		return
	}

	product, _ := h.service.ListProduct(r.Context(), body.ID)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s deleted successfully!!", product.ProductName)
}

// ! UpdateProduct
func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var params repository.UpdateProductParams

	if err := utils.DecodeClient(r, &params); err != nil {
		utils.ErrorHandler(err, w, http.StatusBadRequest)
		return
	}
	if _, err := h.service.UpdateProduct(r.Context(), params); err != nil {
		if err == MinError || err == OneFieldRequired {
			utils.ErrorHandler(err, w, http.StatusBadRequest)
			return
		}

		if err == productNotFoundError {
			utils.ErrorHandler(err, w, http.StatusNotFound)
			return
		}
		utils.ErrorHandler(err, w, http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, utils.SuccessMessage{
		Message: "Product Updated successfully",
	})
}
