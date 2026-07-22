package authutils

import (
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

func PasswordResetValidation(arg PasswordResetParams) error {
	return validate.Struct(arg)
}

func RefreshTokenValidation(arg RefreshTokenParam) error {
	return validate.Struct(arg)
}

func UpdateUsernameValidation(arg UpdateUsernameParams) error {
	return validate.Struct(arg)
}

func UpdateEmailValidation(arg UpdateEmailParams) error {
	return validate.Struct(arg)
}

func UpdateProductValidation(arg repository.UpdateProductParams) error {
	return validate.Struct(UpdateProductParams{
		ProductId: arg.ProductID,
		Price:     arg.Price,
		Quantity:  arg.Quantity,
	})
}
