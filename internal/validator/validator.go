package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
	optsGenValidator "github.com/kazhuravlev/options-gen/pkg/validator"
)

var Validator = validator.New()

func init() {
	optsGenValidator.Set(Validator)

	err := Validator.RegisterValidation("sentrydsn", ValidateSentryDSNOrNil)
	if err != nil {
		panic(err)
	}
}

//nolint:gosimple
func ValidateSentryDSNOrNil(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return true
	}
	newSTR := strings.Split(fl.Field().String(), "://")

	if len(newSTR) != 2 {
		return false
	}

	if newSTR[0] != "http" && newSTR[0] != "https" {
		return false
	}

	newSTR = strings.Split(newSTR[1], "@")
	if len(newSTR) != 2 {
		return false
	}

	return true
}
