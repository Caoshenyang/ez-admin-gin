package response

import (
	"errors"
	"net/http"

	"ez-admin-gin/server/internal/apperror"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Body struct {
	Code    apperror.Code `json:"code"`
	Message string        `json:"message"`
	Data    any           `json:"data,omitempty"`
}

// Success 返回统一成功响应。
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Body{
		Code:    apperror.CodeSuccess,
		Message: "ok",
		Data:    data,
	})
}

// Error 返回统一错误响应。
func Error(c *gin.Context, err error, log *zap.Logger) {
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		c.JSON(appErr.Status, Body{
			Code:    appErr.Code,
			Message: appErr.Message,
		})
		return
	}

	// 未归类错误不把内部细节返回给前端，只记录到日志里。
	if log != nil {
		log.Error("unhandled error", zap.Error(err))
	}

	c.JSON(http.StatusInternalServerError, Body{
		Code:    apperror.CodeInternal,
		Message: "服务器内部错误",
	})
}
