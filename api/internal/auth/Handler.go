package auth

import (
	"encoding/json"
	"net/http"

	repository "github.com/rawbil/ecom2/internal/adapters/sqlc"
	authutils "github.com/rawbil/ecom2/internal/auth/auth-utils"
	"github.com/rawbil/ecom2/internal/utils"
)

type Handler struct {
	Service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		Service: service,
	}
}

// ! REGISTER
func (h *Handler) UserRegister(w http.ResponseWriter, r *http.Request) {
	var registerParams repository.CreateUserParams
	json.NewDecoder(r.Body).Decode(&registerParams)

	_, err := h.Service.UserRegister(r.Context(), registerParams)
	if err != nil {
		if err == FieldsRequiredError {
			utils.ErrorHandler(FieldsRequiredError, w, http.StatusBadRequest)
			return
		}
		if err == InvalidEmailError {
			utils.ErrorHandler(InvalidEmailError, w, http.StatusBadRequest)
			return
		}
		if err == InvalidPasswordError {
			utils.ErrorHandler(InvalidPasswordError, w, http.StatusBadRequest)
			return
		}
		if err == UserExistsError {
			utils.ErrorHandler(UserExistsError, w, http.StatusBadRequest)
			return
		}
		utils.ErrorHandler(err, w, http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, utils.SuccessMessage{
		Message: "Registration success!",
		Data:    nil,
	})
}

// ! LOGIN
func (h *Handler) UserLogin(w http.ResponseWriter, r *http.Request) {
	var params authutils.UserLoginParams
	json.NewDecoder(r.Body).Decode(&params)

	user, token, refreshToken, err := h.Service.UserLogin(r.Context(), params)
	if err != nil {
		if err == InvalidEmailError || err == FieldsRequiredError || err == PasswordMismatchError {
			utils.ErrorHandler(err, w, http.StatusBadRequest)
			return
		}

		if err == UserNotFoundError {
			utils.ErrorHandler(err, w, http.StatusNotFound)
			return
		}

		utils.ErrorHandler(err, w, http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, utils.SuccessMessage{
		Message: "Login Success!",
		Data:    map[string]any{"user": user, "token": token, "refresh_token": refreshToken},
	})
}

// ! LOGOUT
func (h *Handler) UserLogout(w http.ResponseWriter, r *http.Request) {
	err := h.Service.LogoutUser(r.Context())
	if err != nil {
		if err == AuthUserNotFound || err == AuthNotFound {
			utils.ErrorHandler(err, w, http.StatusUnauthorized)
			return
		}
		utils.ErrorHandler(err, w, http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, utils.SuccessMessage{
		Message: "Logged out successfully",
	})
}

// ! PASSWORD RESET
func (h *Handler) PasswordReset(w http.ResponseWriter, r *http.Request) {
	var params authutils.PasswordResetParams
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&params); err != nil {
		utils.ErrorHandler(err, w, http.StatusBadRequest)
		return
	}

	err := h.Service.PasswordReset(r.Context(), params)
	if err != nil {
		if err == FieldsRequiredError || err == InvalidPasswordError {
			utils.ErrorHandler(err, w, http.StatusBadRequest)
			return
		}
		if err == PasswordMismatchError || err == SimilarPasswordError {
			utils.ErrorHandler(err, w, http.StatusBadRequest)
			return
		}
		utils.ErrorHandler(err, w, http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, utils.SuccessMessage{
		Message: "Password reset successful",
	})
}
