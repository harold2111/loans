package utils

import (
	"fmt"
	"github.com/harold2111/loans/shared/errors"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

var Validate *validator.Validate

func InitValidator() {
	Validate = validator.New()
}

func ValidateStruct(s interface{}) error {
	if Validate == nil {
		InitValidator()
	}
	if error := Validate.Struct(s); error != nil {
		firstValidationError := error.(validator.ValidationErrors)[0]

		fmt.Println(firstValidationError.Namespace())
		fmt.Println(firstValidationError.Field())
		fmt.Println(firstValidationError.StructNamespace()) // can differ when a custom TagNameFunc is registered or
		fmt.Println(firstValidationError.StructField())     // by passing alt name to ReportError like below
		fmt.Println(firstValidationError.Tag())
		fmt.Println(firstValidationError.ActualTag())
		fmt.Println(firstValidationError.Kind())
		fmt.Println(firstValidationError.Type())
		fmt.Println(firstValidationError.Value())
		fmt.Println(firstValidationError.Param())

		tag := strings.ToLower(firstValidationError.Tag())
		field := firstValidationError.Field()
		switch tag {
		case "required":
			messagesParameters := []interface{}{field}
			return &errors.ValidationError{ErrorCode: errors.RequiredField, MessagesParameters: messagesParameters}
		default:
			messagesParameters := []interface{}{field}
			return &errors.ValidationError{ErrorCode: errors.InvalidField, MessagesParameters: messagesParameters}
		}

	}
	return nil
}

func ValidateVar(fieldName string, field interface{}, tag string) error {
	if Validate == nil {
		InitValidator()
	}
	if error := Validate.Var(field, tag); error != nil {
		firstValidationError := error.(validator.ValidationErrors)[0]
		fmt.Println(fieldName)
		fmt.Println(firstValidationError.Tag())
		fmt.Println(firstValidationError.ActualTag())
		fmt.Println(firstValidationError.Kind())
		fmt.Println(firstValidationError.Type())
		fmt.Println(firstValidationError.Value())
		fmt.Println(firstValidationError.Param())

		tag := strings.ToLower(firstValidationError.Tag())
		switch tag {
		case "required":
			messagesParameters := []interface{}{fieldName}
			return &errors.ValidationError{ErrorCode: errors.RequiredField, MessagesParameters: messagesParameters}
		default:
			messagesParameters := []interface{}{fieldName}
			return &errors.ValidationError{ErrorCode: errors.InvalidField, MessagesParameters: messagesParameters}
		}

	}
	return nil
}
