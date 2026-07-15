package users

import (
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
		utils.ErrorHandler(err, "oops... server error", w, http.StatusInternalServerError)
		return
	}

	if len(users) < 1 {
		utils.ErrorHandler(errors.New("No user found"), "No user found", w, http.StatusNotFound)
		return
	}

	utils.JsonResponse(w, utils.SuccessMessage{
		Message: "success",
		Data:    map[string]any{"users": users},
	})
}

// ! GET User by email
func (h *Handler) ListUser(w http.ResponseWriter, r *http.Request) {
	type Body struct {
		Email string
	}

	var body Body

	if err := utils.DecodeClient(r, &body); err != nil {
		utils.ErrorHandler(err, "error decoding body", w, http.StatusBadRequest)
		return
	}
	user, err := h.service.ListUser(r.Context(), body.Email)
	if err != nil {
		utils.ErrorHandler(err, "oops... server error", w, http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, utils.SuccessMessage{
		Message: "Success",
		Data:    map[string]any{"user": user},
	})
}

// ! CREATE User
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	params := repository.CreateUserParams{
		Username: "",
		Email:    "",
		Password: "",
	}

	if err := utils.DecodeClient(r, &params); err != nil {
		utils.ErrorHandler(err, err.Error(), w, http.StatusBadRequest)
		return
	}
	if _, err := h.service.CreateUser(r.Context(), params); err != nil {
		utils.ErrorHandler(err, "oops... server error", w, http.StatusInternalServerError)
	}

	utils.JsonResponse(w, utils.SuccessMessage{
		Message: "User Created Successfully",
	})
}

// ! DELETE User
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	type Body struct {
		Email string
	}

	var body Body
	if err := utils.DecodeClient(r, &body); err != nil {
		utils.ErrorHandler(err, "error decoding body", w, http.StatusBadRequest)
		return
	}

	err := h.service.DeleteUser(r.Context(), body.Email)
	if err != nil {
		if err == AuthNotFound {
			utils.ErrorHandler(err, err.Error(), w, http.StatusUnauthorized)
			return
		}

		if err == UserNotFound {
			utils.ErrorHandler(err, err.Error(), w, http.StatusNotFound)
			return
		}

		utils.ErrorHandler(err, "oops... server error", w, http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, utils.SuccessMessage{
		Message: "User Deleted Succesfully",
	})
}
