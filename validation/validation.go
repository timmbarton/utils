package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/timmbarton/utils/types/dates"
	"github.com/timmbarton/utils/types/secs"
)

func New() *validator.Validate {
	v := validator.New()

	err := v.RegisterValidation("dates", DateValidator)
	if err != nil {
		panic(err)
	}

	err = v.RegisterValidation("seconds", SecondsValidator)
	if err != nil {
		panic(err)
	}

	return v
}

func ValidateStruct(val any) error {
	return New().Struct(val)
}

func DateValidator(fl validator.FieldLevel) bool {
	if val, ok := fl.Field().Interface().(dates.Date); ok {
		return val.Unix() > 0
	}

	if val, ok := fl.Field().Interface().(*dates.Date); ok {
		return val != nil && val.Unix() > 0
	}

	return false
}

func SecondsValidator(fl validator.FieldLevel) bool {
	if val, ok := fl.Field().Interface().(secs.Seconds); ok {
		return val.Seconds() > 0
	}

	if val, ok := fl.Field().Interface().(*secs.Seconds); ok {
		return val != nil && val.Seconds() > 0
	}

	return false
}
