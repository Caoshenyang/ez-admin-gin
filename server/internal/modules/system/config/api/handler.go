package api

import (
	"strings"

	configapp "ez-admin-gin/server/internal/modules/system/config/application"
	configdomain "ez-admin-gin/server/internal/modules/system/config/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	httpx "ez-admin-gin/server/internal/pkg/httpx"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	service *configapp.Service
	log     *zap.Logger
}

func NewHandler(service *configapp.Service, log *zap.Logger) *Handler {
	return &Handler{service: service, log: log}
}

func (h *Handler) List(c *gin.Context) {
	var query ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httpx.Error(c, errorsx.BadRequest("查询参数不正确"), h.log)
		return
	}

	result, err := h.service.List(query)
	if err != nil {
		httpx.WriteError(c, err, "查询配置列表失败", h.log)
		return
	}

	httpx.Success(c, result)
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		httpx.WriteError(c, err, "创建系统配置失败", h.log)
		return
	}

	httpx.Success(c, result)
}

func (h *Handler) Update(c *gin.Context) {
	configID, ok := httpx.UintIDParam(c, "id", "配置 ID", h.log)
	if !ok {
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.Update(c.Request.Context(), configID, req)
	if err != nil {
		httpx.WriteError(c, err, "更新系统配置失败", h.log)
		return
	}

	httpx.Success(c, result)
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	configID, ok := httpx.UintIDParam(c, "id", "配置 ID", h.log)
	if !ok {
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	if err := h.service.UpdateStatus(c.Request.Context(), configID, req.Status); err != nil {
		httpx.WriteError(c, err, "更新配置状态失败", h.log)
		return
	}

	httpx.Success(c, gin.H{"id": configID, "status": req.Status})
}

func (h *Handler) Value(c *gin.Context) {
	key := strings.TrimSpace(c.Param("key"))
	result, err := h.service.Value(c.Request.Context(), key)
	if err != nil {
		httpx.WriteError(c, err, "读取系统配置失败", h.log)
		return
	}

	httpx.Success(c, result)
}

var _ = configdomain.PermissionList
