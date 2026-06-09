// Package api 提供接口权限元数据模块的 HTTP 请求处理器与路由定义。
package api

import (
	apiresourceapp "ez-admin-gin/server/internal/modules/iam/apiresource/application"
	httpx "ez-admin-gin/server/internal/pkg/httpx"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	service *apiresourceapp.Service
	log     *zap.Logger
}

func NewHandler(service *apiresourceapp.Service, log *zap.Logger) *Handler {
	return &Handler{service: service, log: log}
}

// List godoc
// @Summary      查询接口权限元数据
// @Tags         IAM / 接口权限
// @Accept       json
// @Produce      json
// @Success      200  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/apis [get]
func (h *Handler) List(c *gin.Context) {
	result, err := h.service.List()
	if err != nil {
		httpx.WriteError(c, err, "查询接口权限元数据失败", h.log)
		return
	}

	httpx.Success(c, result)
}
