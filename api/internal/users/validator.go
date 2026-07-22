package users

import (
	"github.com/go-playground/validator/v10"
	repository "github.com/rawbil/ecom2/internal/adapters/sqlc"
)

type EmailValidation struct {
	Email interface{} `json:"email" validate:"omitempty,email"`
}

func NewValidator() *validator.Validate {
	validator := validator.New()

	return validator
}

var Validator = NewValidator()

func ValidateUserEmail(arg repository.ListUsersParams) error {
	return Validator.Struct(EmailValidation{
		Email: arg.Email,
	})
}