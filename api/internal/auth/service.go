package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-playground/validator/v10"
	repository "github.com/rawbil/ecom2/internal/adapters/sqlc"
	"github.com/rawbil/ecom2/internal/auth/utils"
)

type Service interface {
	UserRegister(ctx context.Context, arg repository.CreateUserParams) (sql.Result, error)

	// UserLogin()
	// UserLogout()
}

type Svc struct {
	repository repository.Queries
}

func NewService(repository repository.Queries) Service {
	return &Svc{
		repository: repository,
	}
}

var (
	FieldsRequiredError  = errors.New("All fields are required")
	InvalidPasswordError = errors.New("Password should have a minimum of 8 characters, have at least 1 uppercase, lowecase letter and special character")
	InvalidEmailError    = errors.New("Invalid email format")
	UserExistsError      = errors.New("User already exists")
	SqlNoRows            = errors.New("No record available")
)

func (svc *Svc) UserRegister(ctx context.Context, params repository.CreateUserParams) (sql.Result, error) {
	//& Validate fields
	if err := utils.UserRegisterValidation(params); err != nil {
		// empty fields
		if ValidationErrorCheck("required", err) {
			return nil, FieldsRequiredError
		}
		// password error
		if ValidationErrorCheck("password_format", err) || ValidationErrorCheck("min", err) {
			return nil, InvalidPasswordError
		}
		// email error
		if ValidationErrorCheck("email", err) {
			return nil, InvalidEmailError
		}
		return nil, err
	}

	//& Ensure user does not exist before registering
	_, err := svc.repository.ListUser(ctx, params.Email)
	if err == nil {
		return nil, UserExistsError
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, UserExistsError
	}

	//& Hash Password
	 hashedPassword, err := utils.PasswordHash(params.Password)
	 if err != nil {
		return nil, err
	}

	params.Password = hashedPassword

	return svc.repository.CreateUser(ctx, params)
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
