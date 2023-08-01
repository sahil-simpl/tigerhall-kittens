package utils

import (
	"time"

	"github.com/go-playground/validator/v10"
)

const (
	dateFormat = "2006-01-02"
)

func ValidateDate(fl validator.FieldLevel) bool {
	_, parseErr := time.Parse(dateFormat, fl.Field().String())

	return parseErr == nil
}
