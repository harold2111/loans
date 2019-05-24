package errors

import "fmt"

type GracefulError struct {
	ErrorCode          string
	MessagesParameters []interface{}
}

func (e *GracefulError) Error() string {
	return errorMessageByErrorCode(e.ErrorCode, e.MessagesParameters)
}

type RecordNotFound struct {
	ErrorCode          string
	MessagesParameters []interface{}
}

func (e *RecordNotFound) Error() string {
	return errorMessageByErrorCode(e.ErrorCode, e.MessagesParameters)
}

type ValidationError struct {
	ErrorCode          string
	MessagesParameters []interface{}
}

func (e *ValidationError) Error() string {
	return errorMessageByErrorCode(e.ErrorCode, e.MessagesParameters)
}

func errorMessageByErrorCode(errorCode string, paramsMessage []interface{}) string {
	if message, ok := ErrorMessages[errorCode]; ok {
		if len(paramsMessage) > 0 {
			v := paramsMessage
			return fmt.Sprintf(message, v...)
		}
		return fmt.Sprintf(message)
	}
	return fmt.Sprintf("Undefined message to errorCode: %v", errorCode)
}
