package custom_error

import "fmt"

// MyError is a custom error type
type CustomError struct {
	Code    int
	Message string
	Origin  string
}

func (ce *CustomError) Error() string {
	return fmt.Sprintf("[%s-%d]: %s", ce.Origin, ce.Code, ce.Message)
}

func FromError(err error, code int, origin string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: err.Error(),
		Origin:  origin,
	}
}
