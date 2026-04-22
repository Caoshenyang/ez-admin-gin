package apperror

import (
	"fmt"
	"net/http"
)

// Code 是业务响应码类型。
// Go 没有 Java 那样的 enum，通常用“自定义类型 + const 常量”表达枚举语义。
type Code int

const (
	// CodeSuccess 表示请求处理成功。
	CodeSuccess Code = 0
	// CodeBadRequest 表示请求参数错误。
	CodeBadRequest Code = 40000
	// CodeUnauthorized 表示未登录或登录已过期。
	CodeUnauthorized Code = 40100
	// CodeForbidden 表示没有权限访问资源。
	CodeForbidden Code = 40300
	// CodeNotFound 表示资源不存在。
	CodeNotFound Code = 40400
	// CodeServiceUnavailable 表示数据库、Redis 等依赖服务不可用。
	CodeServiceUnavailable Code = 50300
	// CodeInternal 表示服务器内部错误。
	CodeInternal Code = 50000
)

// Error 表示可以安全返回给前端的应用错误。
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

// New 创建一个不包裹底层错误的应用错误。
func New(status int, code Code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

// Wrap 创建一个包裹底层错误的应用错误。
func Wrap(err error, status int, code Code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Status:  status,
		Err:     err,
	}
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
