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
	query := r.URL.Query()

	page, _ := strconv.Atoi(query.Get("page"))
	limit, _ := strconv.Atoi(query.Get("limit"))

	product_name := query.Get("name")
	minPrice, _ := strconv.Atoi(query.Get("min_price"))
	maxPrice, _ := strconv.Atoi(query.Get("max_price"))

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	offset := limit * (page - 1)

	products, err := h.service.ListProducts(r.Context(), ListProductsFilter{
		Limit: limit,
		Offset: offset,
		Name: product_name,
		MinPrice: minPrice,
		MaxPrice: maxPrice,
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
