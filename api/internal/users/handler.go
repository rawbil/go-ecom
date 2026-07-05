package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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
	query := r.URL.Query()

	// Pagination
	page := query.Get("page")
	limit := query.Get("limit")

	// Filters
	username := query.Get("username")
	email := query.Get("email")
	role := query.Get("role")

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 || limitInt > 30 {
		limitInt = 10
	}

	offset := (pageInt - 1) * limitInt

	users, err := h.service.ListAllUsers(r.Context(), repository.ListUsersParams{
		Limit:    int32(limitInt),
		Offset:   int32(offset),
		Username: username,
		Email:    email,
		Role:     role,
	})
	if err != nil {
		utils.ErrorHandler(err, w, http.StatusInternalServerError)
		return
	}

	if len(users) < 1 {
		utils.ErrorHandler(errors.New("No user found"), w, http.StatusNotFound)
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

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&body); err != nil {
		utils.ErrorHandler(err, w, http.StatusBadRequest)
		return
	}
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

	if err := utils.DecodeClient(r, &params); err != nil {
		utils.ErrorHandler(err, w, http.StatusBadRequest)
		return
	}
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
