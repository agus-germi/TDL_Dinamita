package api

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	val := validator.New()
	registerCustomValidators(val)
	return val
}

func registerCustomValidators(v *validator.Validate) {
	// Custom validator for datetime format (ISO 8601) using in DTOs
	v.RegisterValidation("datetime", func(fl validator.FieldLevel) bool {
		// Get the format (ISO 8601 for time) from the tag
		format := fl.Param()

		// Get the field value and ensure it is a string
		value, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}

		// Attempt to parse the value with the provided format
		_, err := time.Parse(format, value)
		return err == nil
	})
}
