package validator

import (
	"github.com/go-playground/validator/v10"
)

func GetValidationErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required", "required_if":
		return "should be required"
	case "url":
		return "should be url format"
	case "email":
		return "should be email format"
	case "latitude":
		return "should be latitude format"
	case "longitude":
		return "should be longitude format"
	case "startswith":
		return "should start with " + fe.Param()
	case "e164":
		return "should be phone number (e164) format"
	case "date":
		return "should be date (yyyy-mm-dd) format"
	case "datetime":
		return "should be date (yyyy-mm-dd hh:mm:ss) format"
	case "unique":
		return "should be unique"
	case "min":
		return "length should be more than " + fe.Param() + " characters"
	case "max":
		return "length exceeds " + fe.Param() + " characters"
	case "lte":
		return "should be less than " + fe.Param()
	case "gt":
		return "should be greater than " + fe.Param()
	case "gte":
		return "should be greater or equal than " + fe.Param()
	case "oneof":
		return "should be one of " + fe.Param()
	}

	return fe.Error()
}
