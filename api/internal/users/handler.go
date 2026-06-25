package users

import (
	"encoding/json"
	"net/http"

	repository "github.com/rawbil/ecom2/internal/adapters/sqlc"
	"github.com/rawbil/ecom2/internal/utils"
)

type Handler struct {
	service Service
}

type Response struct {
	Message string
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// ! GET All Users
func (h *Handler) ListAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListAllUsers(r.Context())
	if err != nil {
		utils.ErrorHandler(err, w, http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, users)
}

// ! GET User by email
func (h *Handler) ListUser(w http.ResponseWriter, r *http.Request) {
	type Body struct {
		Email string
	}

	var body Body

	json.NewDecoder(r.Body).Decode(&body)
	user, err := h.service.ListUser(r.Context(), body.Email)
	if err != nil {
		utils.ErrorHandler(err, w, http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, user)
}

// ! CREATE User
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	params := repository.CreateUserParams{
		Username: "",
		Email:    "",
		Password: "",
	}

	json.NewDecoder(r.Body).Decode(&params)
	if _, err := h.service.CreateUser(r.Context(), params); err != nil {
		utils.ErrorHandler(err, w, http.StatusInternalServerError)
	}

	utils.JsonResponse(w, Response{"User Created Successfully!"})
}

// ! DELETE User
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	type Body struct {
		Email string
	}

	var body Body
	json.NewDecoder(r.Body).Decode(&body)

	err := h.service.DeleteUser(r.Context(), body.Email)
	if err != nil {
		utils.ErrorHandler(err, w, http.StatusInternalServerError)
	}
}
