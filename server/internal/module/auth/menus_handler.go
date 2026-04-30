package auth

import (
	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/middleware"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// MenuHandler 负责当前登录用户菜单接口。
type MenuHandler struct {
	service *MenuService
	log     *zap.Logger
}

// NewMenuHandler 创建菜单 Handler。
func NewMenuHandler(service *MenuService, log *zap.Logger) *MenuHandler {
	return &MenuHandler{service: service, log: log}
}

// Menus 返回当前登录用户可见菜单树。
func (h *MenuHandler) Menus(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, apperror.Unauthorized("请先登录"), h.log)
		return
	}

	result, err := h.service.Menus(userID)
	if err != nil {
		writeAuthError(c, err, "查询菜单失败", h.log)
		return
	}

	response.Success(c, result)
}
