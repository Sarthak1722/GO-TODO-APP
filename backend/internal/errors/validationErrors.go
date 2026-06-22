package errors

import "github.com/go-playground/validator/v10"

func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	for _, e := range err.(validator.ValidationErrors) {
		errors[e.Field()] = getErrorMessage(e)
	}

	return errors
}

func getErrorMessage(e validator.FieldError) string {
	switch e.Tag() {

	case "required":
		return "This field is required"

	case "email":
		return "Invalid email format"

	case "min":
		return "Too short"

	case "max":
		return "Too long"

	default:
		return "Invalid field"
	}
}