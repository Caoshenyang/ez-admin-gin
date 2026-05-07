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

// List godoc
// @Summary      查询配置列表
// @Tags         System / 配置管理
// @Accept       json
// @Produce      json
// @Param        page        query     int     false  "页码"
// @Param        page_size   query     int     false  "每页条数"
// @Param        keyword     query     string  false  "关键词"
// @Param        group_code  query     string  false  "配置分组"
// @Param        status      query     int     false  "状态"
// @Success      200  {object}  httpx.Body{data=configdomain.ListResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/configs [get]
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

// Create godoc
// @Summary      创建系统配置
// @Tags         System / 配置管理
// @Accept       json
// @Produce      json
// @Param        body  body  configdomain.CreateRequest  true  "配置参数"
// @Success      200  {object}  httpx.Body{data=configdomain.Response}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/configs [post]
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

// Update godoc
// @Summary      更新系统配置
// @Tags         System / 配置管理
// @Accept       json
// @Produce      json
// @Param        id    path  uint                        true  "配置 ID"
// @Param        body  body  configdomain.UpdateRequest  true  "配置参数"
// @Success      200  {object}  httpx.Body{data=configdomain.Response}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/configs/{id}/update [post]
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

// UpdateStatus godoc
// @Summary      更新配置状态
// @Tags         System / 配置管理
// @Accept       json
// @Produce      json
// @Param        id    path  uint                              true  "配置 ID"
// @Param        body  body  configdomain.UpdateStatusRequest  true  "状态参数"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/configs/{id}/status [post]
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

// Value godoc
// @Summary      读取配置值
// @Tags         System / 配置管理
// @Accept       json
// @Produce      json
// @Param        key  path  string  true  "配置键"
// @Success      200  {object}  httpx.Body{data=configdomain.ValueResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/configs/value/{key} [get]
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
