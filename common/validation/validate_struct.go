package validation

import (
	"reflect"

	"github.com/go-playground/validator/v10"

	"github.com/danielkhtse/supreme-adventure/common/types"
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

	v.RegisterValidation("transaction_status", validateTransactionStatusEnum)

	return nil
}

func validateTransactionStatusEnum(fl validator.FieldLevel) bool {
	transactionStatus, ok := fl.Field().Interface().(types.TransactionStatus)
	if !ok {
		return false
	}

	switch transactionStatus {
	case types.TransactionStatusPending,
		types.TransactionStatusCompleted,
		types.TransactionStatusFailed:
		return true
	default:
		return false
	}
}
