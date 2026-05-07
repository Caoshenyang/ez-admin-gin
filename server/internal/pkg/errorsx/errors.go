package errorsx

import (
	"fmt"
	"net/http"
)

type Code int

const (
	CodeSuccess            Code = 0
	CodeBadRequest         Code = 40000
	CodeUnauthorized       Code = 40100
	CodeForbidden          Code = 40300
	CodeNotFound           Code = 40400
	CodeServiceUnavailable Code = 50300
	CodeInternal           Code = 50000
)

type Error struct {
	Code    Code
	Message string
	Status  int
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}

	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func New(status int, code Code, message string) *Error {
	return &Error{Code: code, Message: message, Status: status}
}

func Wrap(err error, status int, code Code, message string) *Error {
	return &Error{Code: code, Message: message, Status: status, Err: err}
}

func BadRequest(message string) *Error {
	return New(http.StatusBadRequest, CodeBadRequest, message)
}

func Unauthorized(message string) *Error {
	return New(http.StatusUnauthorized, CodeUnauthorized, message)
}

func Forbidden(message string) *Error {
	return New(http.StatusForbidden, CodeForbidden, message)
}

func NotFound(message string) *Error {
	return New(http.StatusNotFound, CodeNotFound, message)
}

func ServiceUnavailable(message string, err error) *Error {
	return Wrap(err, http.StatusServiceUnavailable, CodeServiceUnavailable, message)
}

func Internal(message string, err error) *Error {
	return Wrap(err, http.StatusInternalServerError, CodeInternal, message)
}
