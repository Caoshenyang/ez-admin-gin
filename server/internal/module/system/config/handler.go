package config

import (
	"errors"
	"strconv"
	"strings"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 负责系统配置模块的 HTTP 协议层绑定与输出。
type Handler struct {
	service *Service
	log     *zap.Logger
}

// NewHandler 创建配置 Handler。
func NewHandler(service *Service, log *zap.Logger) *Handler {
	return &Handler{service: service, log: log}
}

// List 返回系统配置分页列表。
func (h *Handler) List(c *gin.Context) {
	var query ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, apperror.BadRequest("查询参数不正确"), h.log)
		return
	}

	result, err := h.service.List(query)
	if err != nil {
		writeError(c, err, "查询配置列表失败", h.log)
		return
	}

	response.Success(c, result)
}

// Create 创建系统配置。
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		writeError(c, err, "创建系统配置失败", h.log)
		return
	}

	response.Success(c, result)
}

// Update 编辑系统配置。
func (h *Handler) Update(c *gin.Context) {
	configID, ok := configIDParam(c, h.log)
	if !ok {
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.Update(c.Request.Context(), configID, req)
	if err != nil {
		writeError(c, err, "更新系统配置失败", h.log)
		return
	}

	response.Success(c, result)
}

// UpdateStatus 单独修改配置状态。
func (h *Handler) UpdateStatus(c *gin.Context) {
	configID, ok := configIDParam(c, h.log)
	if !ok {
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	if err := h.service.UpdateStatus(c.Request.Context(), configID, req.Status); err != nil {
		writeError(c, err, "更新配置状态失败", h.log)
		return
	}

	response.Success(c, gin.H{"id": configID, "status": req.Status})
}

// Value 按配置键读取启用中的配置值。
func (h *Handler) Value(c *gin.Context) {
	key := strings.TrimSpace(c.Param("key"))
	result, err := h.service.Value(c.Request.Context(), key)
	if err != nil {
		writeError(c, err, "读取系统配置失败", h.log)
		return
	}

	response.Success(c, result)
}

func configIDParam(c *gin.Context, log *zap.Logger) (uint, bool) {
	rawID := c.Param("id")
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 {
		response.Error(c, apperror.BadRequest("配置 ID 不正确"), log)
		return 0, false
	}

	return uint(id), true
}

func writeError(c *gin.Context, err error, fallbackMessage string, log *zap.Logger) {
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		response.Error(c, appErr, log)
		return
	}

	response.Error(c, apperror.Internal(fallbackMessage, err), log)
}
