package utils

import "github.com/go-playground/validator/v10"

type UserRegisterParams struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required, email"`
	Password string `json:"password" validate:"required, password_format"`
}

// var Validate = validator.New()