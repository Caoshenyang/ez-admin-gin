package notice

import (
	"errors"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 负责公告模块的 HTTP 协议层绑定与输出。
type Handler struct {
	service *Service
	log     *zap.Logger
}

// NewHandler 创建公告 Handler。
func NewHandler(service *Service, log *zap.Logger) *Handler {
	return &Handler{service: service, log: log}
}

// List 返回公告分页列表。
func (h *Handler) List(c *gin.Context) {
	var query ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, apperror.BadRequest("查询参数不正确"), h.log)
		return
	}

	result, err := h.service.List(query)
	if err != nil {
		writeError(c, err, "查询公告列表失败", h.log)
		return
	}

	response.Success(c, result)
}

// Create 创建公告。
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.Create(req)
	if err != nil {
		writeError(c, err, "创建公告失败", h.log)
		return
	}

	response.Success(c, result)
}

// Update 编辑公告。
func (h *Handler) Update(c *gin.Context) {
	noticeID, err := ParseNoticeID(c.Param("id"))
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.Update(noticeID, req)
	if err != nil {
		writeError(c, err, "更新公告失败", h.log)
		return
	}

	response.Success(c, result)
}

// UpdateStatus 修改公告状态。
func (h *Handler) UpdateStatus(c *gin.Context) {
	noticeID, err := ParseNoticeID(c.Param("id"))
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	if err := h.service.UpdateStatus(noticeID, req.Status); err != nil {
		writeError(c, err, "更新公告状态失败", h.log)
		return
	}

	response.Success(c, gin.H{
		"id":     noticeID,
		"status": req.Status,
	})
}

func writeError(c *gin.Context, err error, fallbackMessage string, log *zap.Logger) {
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		response.Error(c, appErr, log)
		return
	}

	response.Error(c, apperror.Internal(fallbackMessage, err), log)
}
