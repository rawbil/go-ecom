package auth

import (
	"encoding/json"
	"net/http"

	repository "github.com/rawbil/ecom2/internal/adapters/sqlc"
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
