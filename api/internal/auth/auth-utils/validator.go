package authutils

import (
	"errors"
	"regexp"

	"github.com/go-playground/validator/v10"
	repository "github.com/rawbil/ecom2/internal/adapters/sqlc"
)

var validate = NewValidator()

func NewValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("password_format", ValidatePasswordFormat)
	return validate
}

var (
	hasUpper       = regexp.MustCompile(`[A-Z]`)
	hasLower       = regexp.MustCompile(`[a-z]`)
	hasSpecialChar = regexp.MustCompile(`[^A-Za-z0-9]`)
)

func ValidatePasswordFormat(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return hasUpper.MatchString(password) &&
		hasLower.MatchString(password) &&
		hasSpecialChar.MatchString(password)
}

func UserRegisterValidation(args repository.CreateUserParams) error {
	return validate.Struct(UserRegisterParams{
		Username: args.Username,
		Email:    args.Email,
		Password: args.Password,
	})
}

func UserLoginValidation(arg UserLoginParams) error {
	return validate.Struct(arg)
}

func ValidationErrorCheck(tag string, err error) bool {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return false
	}

	for _, ValidationError := range validationErrors {
		if ValidationError.Tag() == tag {
			return true
		}
	}

	return false
}