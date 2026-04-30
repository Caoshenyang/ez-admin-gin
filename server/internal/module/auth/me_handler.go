package auth

import (
	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/middleware"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// MeHandler 负责当前登录用户接口。
type MeHandler struct {
	service *MeService
	log     *zap.Logger
}

// NewMeHandler 创建当前用户 Handler。
func NewMeHandler(service *MeService, log *zap.Logger) *MeHandler {
	return &MeHandler{service: service, log: log}
}

// Me 返回当前登录用户基础信息与数据范围摘要。
func (h *MeHandler) Me(c *gin.Context) {
	if actor, ok := middleware.CurrentActor(c); ok {
		response.Success(c, h.service.Build(actor))
		return
	}

	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, apperror.Unauthorized("请先登录"), h.log)
		return
	}

	username, _ := middleware.CurrentUsername(c)
	response.Success(c, MeResponse{
		UserID:   userID,
		Username: username,
	})
}
