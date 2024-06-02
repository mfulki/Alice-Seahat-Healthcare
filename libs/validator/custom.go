package validator

import (
	"reflect"
	"strings"
	"time"

	"Alice-Seahat-Healthcare/seahat-be/constant"

	"github.com/go-playground/validator/v10"
)

func SetCustom(validate any) {
	if v, ok := validate.(*validator.Validate); ok {
		v.RegisterTagNameFunc(fieldTagNew)
		v.RegisterValidation("date", isDateTimeFormat(constant.DateFormat))
		v.RegisterValidation("datetime", isDateTimeFormat(constant.FullTimeFormat))
	}
}

func fieldTag(field reflect.StructField, tagName string) string {
	name := strings.SplitN(field.Tag.Get(tagName), ",", 2)[0]
	if name == "-" {
		return ""
	}

	return name
}

func fieldTagNew(field reflect.StructField) string {
	name := fieldTag(field, "json")
	if name == "" {
		name = fieldTag(field, "form")
	}

	return name
}

func isDateTimeFormat(format string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		_, err := time.Parse(format, value)

		return err == nil
	}
}
