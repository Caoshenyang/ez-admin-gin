package auth

import (
	"errors"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoginHandler 负责登录接口的协议层绑定与输出。
type LoginHandler struct {
	service *LoginService
	log     *zap.Logger
}

// NewLoginHandler 创建登录 Handler。
func NewLoginHandler(service *LoginService, log *zap.Logger) *LoginHandler {
	return &LoginHandler{service: service, log: log}
}

// Login 校验用户名和密码并返回登录态。
func (h *LoginHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.service.recordLogin(c.Request.Context(), 0, "", 2, "用户名和密码不能为空", c.ClientIP(), c.Request.UserAgent())
		response.Error(c, apperror.BadRequest("用户名和密码不能为空"), h.log)
		return
	}

	result, err := h.service.Login(c.Request.Context(), req, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		writeAuthError(c, err, "登录失败", h.log)
		return
	}

	response.Success(c, result)
}

func writeAuthError(c *gin.Context, err error, fallbackMessage string, log *zap.Logger) {
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		response.Error(c, appErr, log)
		return
	}

	response.Error(c, apperror.Internal(fallbackMessage, err), log)
}
