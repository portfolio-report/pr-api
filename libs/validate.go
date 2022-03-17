package libs

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func isDateYYYYMMDD(fl validator.FieldLevel) bool {
	DateYYYYMMDDRegex := regexp.MustCompile(`^[0-9]{4}-(1[0-2]|0[1-9])-(3[01]|[012][0-9])$`)
	return DateYYYYMMDDRegex.MatchString(fl.Field().String())
}

func isLaxUuid(fl validator.FieldLevel) bool {
	LaxUuidRegex := regexp.MustCompile(`^([a-f0-9]{32}|[a-f0-9]{8}-([a-f0-9]{4}-){3}[a-f0-9]{12})$`)
	return LaxUuidRegex.MatchString(fl.Field().String())
}

func isValidUsername(fl validator.FieldLevel) bool {
	ValidUserNameRegex := regexp.MustCompile(`^([a-zA-Z0-9._-])*$`)
	return ValidUserNameRegex.MatchString(fl.Field().String())
}

func isValidDecimalString(fl validator.FieldLevel) bool {
	ValidDecimalStringRegex := regexp.MustCompile(`^[0-9]+\.?[0-9]*$`)
	return ValidDecimalStringRegex.MatchString(fl.Field().String())
}

// RegisterCustomValidations registers custom validators
func RegisterCustomValidations(v *validator.Validate) {
	v.RegisterValidation("DateYYYY-MM-DD", isDateYYYYMMDD)
	v.RegisterValidation("LaxUuid", isLaxUuid)
	v.RegisterValidation("ValidUsername", isValidUsername)
	v.RegisterValidation("DecimalString", isValidDecimalString)

}

// GetValidator creates and returns validator with custom validations
func GetValidator() *validator.Validate {
	v := validator.New()
	RegisterCustomValidations(v)
	return v
}
