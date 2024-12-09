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
	v.RegisterValidation("datetime", func(fl validator.FieldLevel) bool {
		format := fl.Param()

		value, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		_, err := time.Parse(format, value)
		return err == nil
	})
}
