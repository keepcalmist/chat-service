package errors

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
)

// ServerError is used to return custom error codes to client.
type ServerError struct {
	Code    int
	Message string
	cause   error
}

func NewServerError(code int, msg string, err error) *ServerError {
	return &ServerError{
		Code:    code,
		Message: msg,
		cause:   err,
	}
}

func (s *ServerError) Error() string {
	return fmt.Sprintf("%s: %v", s.Message, s.cause)
}

func (s *ServerError) Is(err error) bool {
	return errors.Is(s.cause, err)
}

func GetServerErrorCode(err error) int {
	code, _, _ := ProcessServerError(err)
	return code
}

// ProcessServerError tries to retrieve from given error it's code, message and some details.
// For example, that fields can be used to build error response for client.
func ProcessServerError(err error) (code int, msg string, details string) {
	if serverError := new(ServerError); errors.As(err, &serverError) {
		return serverError.Code, serverError.Message, serverError.Error()
	}

	if echoErr := new(echo.HTTPError); errors.As(err, &echoErr) {
		return echoErr.Code, echoErr.Message.(string), echoErr.Error()
	}

	return 500, "something went wrong", err.Error()
}
