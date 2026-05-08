// Package errorsx 提供带 HTTP 状态码和业务错误码的应用层错误类型。
package errorsx

import (
	"fmt"
	"net/http"
)

// Code 是业务错误码类型。
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

// Error 是应用层统一错误类型，携带 HTTP 状态码、业务错误码和用户可见消息。
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

// New 创建一个不包含原始错误的 Error。
func New(status int, code Code, message string) *Error {
	return &Error{Code: code, Message: message, Status: status}
}

// Wrap 创建一个包含原始错误的 Error。
func Wrap(err error, status int, code Code, message string) *Error {
	return &Error{Code: code, Message: message, Status: status, Err: err}
}

// BadRequest 返回 400 Bad Request 错误。
func BadRequest(message string) *Error {
	return New(http.StatusBadRequest, CodeBadRequest, message)
}

// Unauthorized 返回 401 Unauthorized 错误。
func Unauthorized(message string) *Error {
	return New(http.StatusUnauthorized, CodeUnauthorized, message)
}

// Forbidden 返回 403 Forbidden 错误。
func Forbidden(message string) *Error {
	return New(http.StatusForbidden, CodeForbidden, message)
}

// NotFound 返回 404 Not Found 错误。
func NotFound(message string) *Error {
	return New(http.StatusNotFound, CodeNotFound, message)
}

// ServiceUnavailable 返回 503 Service Unavailable 错误。
func ServiceUnavailable(message string, err error) *Error {
	return Wrap(err, http.StatusServiceUnavailable, CodeServiceUnavailable, message)
}

// Internal 返回 500 Internal Server Error 错误。
func Internal(message string, err error) *Error {
	return Wrap(err, http.StatusInternalServerError, CodeInternal, message)
}
