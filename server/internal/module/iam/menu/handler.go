package menu

import (
	"errors"
	"strconv"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 负责菜单模块的 HTTP 协议层绑定与输出。
type Handler struct {
	service *Service
	log     *zap.Logger
}

// NewHandler 创建菜单 Handler。
func NewHandler(service *Service, log *zap.Logger) *Handler {
	return &Handler{service: service, log: log}
}

// List 返回完整菜单树。
func (h *Handler) List(c *gin.Context) {
	result, err := h.service.List()
	if err != nil {
		writeError(c, err, "查询菜单树失败", h.log)
		return
	}

	response.Success(c, result)
}

// Create 创建菜单、目录或按钮。
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.Create(req)
	if err != nil {
		writeError(c, err, "创建菜单失败", h.log)
		return
	}

	response.Success(c, result)
}

// Update 编辑菜单基础信息。
func (h *Handler) Update(c *gin.Context) {
	menuID, ok := menuIDParam(c, h.log)
	if !ok {
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.Update(menuID, req)
	if err != nil {
		writeError(c, err, "更新菜单失败", h.log)
		return
	}

	response.Success(c, result)
}

// UpdateStatus 单独更新菜单状态。
func (h *Handler) UpdateStatus(c *gin.Context) {
	menuID, ok := menuIDParam(c, h.log)
	if !ok {
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	if err := h.service.UpdateStatus(menuID, req.Status); err != nil {
		writeError(c, err, "更新菜单状态失败", h.log)
		return
	}

	response.Success(c, gin.H{"id": menuID, "status": req.Status})
}

// Delete 删除菜单节点。
func (h *Handler) Delete(c *gin.Context) {
	menuID, ok := menuIDParam(c, h.log)
	if !ok {
		return
	}

	if err := h.service.Delete(menuID); err != nil {
		writeError(c, err, "删除菜单失败", h.log)
		return
	}

	response.Success(c, gin.H{"id": menuID})
}

func menuIDParam(c *gin.Context, log *zap.Logger) (uint, bool) {
	rawID := c.Param("id")
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 {
		response.Error(c, apperror.BadRequest("菜单 ID 不正确"), log)
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
