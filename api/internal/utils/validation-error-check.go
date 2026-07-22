package utils

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

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
