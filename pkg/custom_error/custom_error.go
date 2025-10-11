package custom_error

import "fmt"

// MyError is a custom error type
type CustomError struct {
	Code    int
	Message string
}

func (ce *CustomError) Error() string {
	return fmt.Sprintf("Code %d: %s", ce.Code, ce.Message)
}

func FromError(err error, code int) *CustomError {
	return &CustomError{
		Code:    code,
		Message: err.Error(),
	}
}
