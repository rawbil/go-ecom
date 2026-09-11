package auth

import (
	"encoding/json"
	"net/http"

	repository "github.com/rawbil/ecom2/internal/adapters/sqlc"
	authutils "github.com/rawbil/ecom2/internal/auth/auth-utils"
	"github.com/rawbil/ecom2/internal/config"
	"github.com/rawbil/ecom2/internal/utils"
)

type Handler struct {
	Service Service
}

type RegisterEmailData struct {
	Username string
}

func NewHandler(service Service) *Handler {
	return &Handler{
		Service: service,
	}
}

// ! REGISTER
func (h *Handler) UserRegister(w http.ResponseWriter, r *http.Request) {
	var registerParams repository.CreateUserParams

	if err := utils.DecodeClient(r, &registerParams); err != nil {
		utils.ErrorHandler(err, "error decoding body", w, http.StatusBadRequest)
		return
	}

	_, err := h.Service.UserRegister(r.Context(), registerParams)
	if err != nil {
		if err == FieldsRequiredError || err == InvalidEmailError || err == InvalidPasswordError || err == UsernameLenErr {
			utils.ErrorHandler(FieldsRequiredError, err.Error(), w, http.StatusBadRequest)
			return
		}
		if err == UserExistsError {
			utils.ErrorHandler(UserExistsError, err.Error(), w, http.StatusConflict)
			return
		}
		utils.ErrorHandler(err, "oops... server error", w, http.StatusInternalServerError)
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
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&params); err != nil {
		utils.ErrorHandler(err, "error decoding body", w, http.StatusBadRequest)
		return
	}

	user, token, refreshToken, err := h.Service.UserLogin(r.Context(), params)
	if err != nil {
		if err == InvalidEmailError || err == FieldsRequiredError || err == PasswordMismatchError {
			utils.ErrorHandler(err, err.Error(), w, http.StatusBadRequest)
			return
		}

		if err == UserNotFoundError {
			utils.ErrorHandler(err, err.Error(), w, http.StatusNotFound)
			return
		}

		utils.ErrorHandler(err, "oops... server error", w, http.StatusInternalServerError)
		return
	}

	rt_cookie := http.Cookie{
		Name:     "__refresh_token__",
		Value:    refreshToken,
		MaxAge:   3600, // 1 hour in seconds
		Path:     "/",
		HttpOnly: true,                                         // Prevents client-side JS access
		Secure:   config.GetServerConfigFunc().APPENV != "dev", // true in prod (SET APP_ENV=prod)
	}

	http.SetCookie(w, &rt_cookie)

	utils.JsonResponse(w, utils.SuccessMessage{
		Message: "Login Success!",
		Data:    map[string]any{"user": user, "access_token": token},
	})
}

// ! LOGOUT
func (h *Handler) UserLogout(w http.ResponseWriter, r *http.Request) {
	err := h.Service.LogoutUser(r.Context())
	if err != nil {
		if err == AuthUserNotFound || err == AuthNotFound {
			utils.ErrorHandler(err, err.Error(), w, http.StatusUnauthorized)
			return
		}
		utils.ErrorHandler(err, "oops... server error", w, http.StatusInternalServerError)
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
		utils.ErrorHandler(err, "error decoding body", w, http.StatusBadRequest)
		return
	}

	err := h.Service.PasswordReset(r.Context(), params)
	if err != nil {
		if err == FieldsRequiredError || err == InvalidPasswordError {
			utils.ErrorHandler(err, err.Error(), w, http.StatusBadRequest)
			return
		}
		if err == PasswordMismatchError || err == SimilarPasswordError {
			utils.ErrorHandler(err, err.Error(), w, http.StatusBadRequest)
			return
		}
		utils.ErrorHandler(err, "oops... server error", w, http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, utils.SuccessMessage{
		Message: "Password reset successful",
	})
}

// ! REFRESH TOKENS
func (h *Handler) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	var param authutils.RefreshTokenParam

	if err := utils.DecodeClient(r, &param); err != nil {
		utils.ErrorHandler(err, "error decoding body", w, http.StatusBadRequest)
		return
	}

	authToken, refreshToken, err := h.Service.RefreshTokens(r.Context(), param)
	if err != nil {
		if err == AuthNotFound || err == AuthUserNotFound || err == InvalidRefreshToken || err == TokenExpiredError {
			utils.ErrorHandler(err, err.Error(), w, http.StatusUnauthorized)
			return
		}

		if err == FieldsRequiredError {
			utils.ErrorHandler(err, err.Error(), w, http.StatusBadRequest)
			return
		}
		utils.ErrorHandler(err, "oops... server error", w, http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, utils.SuccessMessage{
		Message: "Success",
		Data:    map[string]string{"access_token": authToken, "refresh_token": refreshToken},
	})
}

// ! UPDATE MY USERNAME
func (h *Handler) UpdateUsername(w http.ResponseWriter, r *http.Request) {
	var params authutils.UpdateUsernameParams

	if err := utils.DecodeClient(r, &params); err != nil {
		utils.ErrorHandler(err, "error decoding body", w, http.StatusBadRequest)
		return
	}

	user, err := h.Service.UpdateMyUsername(r.Context(), params)
	if err != nil {

		if err == AuthUserNotFound || err == UserNotFoundError {
			utils.ErrorHandler(err, err.Error(), w, http.StatusUnauthorized)
			return
		}

		if err == FieldsRequiredError || err == UsernameLenErr {
			utils.ErrorHandler(err, err.Error(), w, http.StatusBadRequest)
			return
		}

		utils.ErrorHandler(err, "oops... server error", w, http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, utils.SuccessMessage{
		Message: "details updated",
		Data:    map[string]any{"user": user},
	})
}

// ! UPDATE MY EMAIL
func (h *Handler) UpdateUserEmail(w http.ResponseWriter, r *http.Request) {
	var params authutils.UpdateEmailParams

	if err := utils.DecodeClient(r, &params); err != nil {
		utils.ErrorHandler(err, "error decoding body", w, http.StatusBadRequest)
		return
	}

	user, err := h.Service.UpdateMyEmail(r.Context(), params)
	if err != nil {

		if err == AuthUserNotFound || err == UserNotFoundError {
			utils.ErrorHandler(err, err.Error(), w, http.StatusUnauthorized)
			return
		}

		if err == FieldsRequiredError || err == InvalidEmailError {
			utils.ErrorHandler(err, err.Error(), w, http.StatusBadRequest)
			return
		}

		if err == EmailTaken {
			utils.ErrorHandler(err, err.Error(), w, http.StatusConflict)
			return
		}

		utils.ErrorHandler(err, "oops... server error", w, http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, utils.SuccessMessage{
		Message: "details updated",
		Data:    map[string]any{"user": user},
	})
}

// ! FORGOT PASSWORD
func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var params authutils.UpdateEmailParams

	if err := utils.DecodeClient(r, &params); err != nil {
		utils.ErrorHandler(err, "error decoding body", w, http.StatusBadRequest)
		return
	}

	err := h.Service.ForgotPassword(r.Context(), params)
	if err != nil {
		if err == AuthUserNotFound || err == UserNotFoundError {
			utils.ErrorHandler(err, err.Error(), w, http.StatusUnauthorized)
			return
		}

		if err == FieldsRequiredError || err == InvalidEmailError {
			utils.ErrorHandler(err, err.Error(), w, http.StatusBadRequest)
			return
		}
		utils.ErrorHandler(err, err.Error(), w, http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, utils.SuccessMessage{
		Message: "Success",
	})
}
