package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidationError formats validation errors into field-based error slices
type ValidationError map[string][]string

// Validate runs validation tags on a struct and returns formatted error messages
func Validate(data interface{}) ValidationError {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return ValidationError{"global": []string{err.Error()}}
	}

	errs := make(ValidationError)
	for _, e := range validationErrors {
		field := strings.ToLower(e.Field())
		var message string

		switch e.Tag() {
		case "required":
			message = fmt.Sprintf("Kolom %s wajib diisi.", field)
		case "email":
			message = "Format email tidak valid."
		case "min":
			message = fmt.Sprintf("Kolom %s minimal %s karakter.", field, e.Param())
		case "max":
			message = fmt.Sprintf("Kolom %s maksimal %s karakter.", field, e.Param())
		case "eqfield":
			message = fmt.Sprintf("Kolom %s tidak cocok.", field)
		case "oneof":
			message = fmt.Sprintf("Kolom %s harus salah satu dari: %s.", field, e.Param())
		default:
			message = fmt.Sprintf("Kolom %s tidak valid (%s).", field, e.Tag())
		}

		errs[field] = append(errs[field], message)
	}

	return errs
}
