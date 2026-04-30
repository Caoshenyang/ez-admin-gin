package auth

import (
	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/middleware"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// DashboardHandler 负责工作台概览接口。
type DashboardHandler struct {
	service *DashboardService
	log     *zap.Logger
}

// NewDashboardHandler 创建工作台 Handler。
func NewDashboardHandler(service *DashboardService, log *zap.Logger) *DashboardHandler {
	return &DashboardHandler{service: service, log: log}
}

// Dashboard 返回工作台概览数据。
func (h *DashboardHandler) Dashboard(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, apperror.Unauthorized("请先登录"), h.log)
		return
	}

	username, _ := middleware.CurrentUsername(c)
	result, err := h.service.Dashboard(userID, username)
	if err != nil {
		writeAuthError(c, err, "查询工作台失败", h.log)
		return
	}

	response.Success(c, result)
}
