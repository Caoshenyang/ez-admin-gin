// Package httpx 封装统一的 JSON 响应写入和错误处理。
package httpx

import (
	"errors"
	"net/http"
	"strconv"

	"ez-admin-gin/server/internal/pkg/actorx"
	"ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/datascope"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Body 是统一的 JSON 响应体结构。
type Body struct {
	Code    errorsx.Code `json:"code"`
	Message string       `json:"message"`
	Data    any          `json:"data,omitempty"`
}

// Success 写入成功的 JSON 响应。
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Body{
		Code:    errorsx.CodeSuccess,
		Message: "ok",
		Data:    data,
	})
}

// Error 写入错误 JSON 响应，自动识别 errorsx.Error 并设置对应 HTTP 状态码。
func Error(c *gin.Context, err error, log *zap.Logger) {
	var appErr *errorsx.Error
	if errors.As(err, &appErr) {
		c.JSON(appErr.Status, Body{
			Code:    appErr.Code,
			Message: appErr.Message,
		})
		return
	}

	if log != nil {
		log.Error("unhandled error", zap.Error(err))
	}

	c.JSON(http.StatusInternalServerError, Body{
		Code:    errorsx.CodeInternal,
		Message: "服务器内部错误",
	})
}

// WriteError 写入错误 JSON 响应，非 errorsx.Error 时用 fallbackMessage 包装为内部错误。
func WriteError(c *gin.Context, err error, fallbackMessage string, log *zap.Logger) {
	var appErr *errorsx.Error
	if errors.As(err, &appErr) {
		Error(c, appErr, log)
		return
	}

	Error(c, errorsx.Internal(fallbackMessage, err), log)
}

// CurrentActor 从 Gin 上下文中读取当前登录人信息，未登录时自动写入 401 响应。
func CurrentActor(c *gin.Context, log *zap.Logger) (datascope.Actor, bool) {
	actor, ok := actorx.CurrentActor(c)
	if !ok {
		Error(c, errorsx.Unauthorized("请先登录"), log)
		return datascope.Actor{}, false
	}
	return actor, true
}

// UintIDParam 从 URL 路径参数中解析 uint 类型的 ID，无效时自动写入 400 响应。
func UintIDParam(c *gin.Context, param string, label string, log *zap.Logger) (uint, bool) {
	rawID := c.Param(param)
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 {
		Error(c, errorsx.BadRequest(label+" 不正确"), log)
		return 0, false
	}

	return uint(id), true
}
