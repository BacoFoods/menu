package shared

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var OneOfValidationTranslator = func(ut ut.Translator) error {
	return ut.Add("oneof", "the field {0} must be one of {1}", true)
}

var OneOfValidation = func(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T("oneof", fe.Field(), fe.Param())
	return t
}

var RequiredValidationTranslator = func(ut ut.Translator) error {
	return ut.Add("required", "the field {0} is required", true)
}

var RequiredValidation = func(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T("required", fe.Field())
	return t
}

var HourValidationTranslator = func(ut ut.Translator) error {
	return ut.Add("availablehour", "the field {0} must be a valid 24 hours format like 14:30", true)
}

var HourValidation = func(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T("availablehour", fe.Field())
	return t
}
