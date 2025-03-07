package validation

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

func ValidateStruct(reqPayload interface{}) error {
	v := validator.New()

	v.RegisterValidation("sqlnullint64", func(fl validator.FieldLevel) bool {
		field := fl.Field()
		if field.Kind() != reflect.Struct {
			return false
		}

		intField := field.FieldByName("Int64")
		validField := field.FieldByName("Valid")

		if !validField.Bool() {
			return true
		}

		return intField.Int() <= 9223372036854775807
	})

	v.RegisterValidation("sqlnullfloat64", func(fl validator.FieldLevel) bool {
		field := fl.Field()
		if field.Kind() != reflect.Struct {
			return false
		}

		// floatField := field.FieldByName("Float64")
		validField := field.FieldByName("Valid")

		if !validField.Bool() {
			return true
		}

		return true
	})

	return nil
}
