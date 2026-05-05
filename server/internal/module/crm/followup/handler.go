package followup

import (
	"errors"
	"strconv"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/middleware"
	"ez-admin-gin/server/internal/platform/datascope"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 负责 CRM 客户跟进模块的 HTTP 协议层绑定与输出。
type Handler struct {
	service *Service
	log     *zap.Logger
}

// NewHandler 创建 CRM 客户跟进 Handler。
func NewHandler(service *Service, log *zap.Logger) *Handler {
	return &Handler{service: service, log: log}
}

// ListCustomerOptions 返回当前数据范围内可选客户。
func (h *Handler) ListCustomerOptions(c *gin.Context) {
	actor, ok := currentActor(c, h.log)
	if !ok {
		return
	}

	limit := 0
	if rawLimit := c.Query("limit"); rawLimit != "" {
		parsed, err := strconv.Atoi(rawLimit)
		if err != nil {
			response.Error(c, apperror.BadRequest("limit 参数不正确"), h.log)
			return
		}
		limit = parsed
	}

	result, err := h.service.ListCustomerOptions(actor, c.Query("keyword"), limit)
	if err != nil {
		writeError(c, err, "查询客户选项失败", h.log)
		return
	}

	response.Success(c, result)
}

// List 返回客户跟进分页列表。
func (h *Handler) List(c *gin.Context) {
	actor, ok := currentActor(c, h.log)
	if !ok {
		return
	}

	var query ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, apperror.BadRequest("查询参数不正确"), h.log)
		return
	}

	result, err := h.service.List(actor, query)
	if err != nil {
		writeError(c, err, "查询客户跟进列表失败", h.log)
		return
	}

	response.Success(c, result)
}

// Create 创建客户跟进。
func (h *Handler) Create(c *gin.Context) {
	actor, ok := currentActor(c, h.log)
	if !ok {
		return
	}

	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.Create(actor, req)
	if err != nil {
		writeError(c, err, "创建客户跟进失败", h.log)
		return
	}

	response.Success(c, result)
}

// Update 编辑客户跟进。
func (h *Handler) Update(c *gin.Context) {
	actor, ok := currentActor(c, h.log)
	if !ok {
		return
	}

	followUpID, err := ParseFollowUpID(c.Param("id"))
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.Update(actor, followUpID, req)
	if err != nil {
		writeError(c, err, "更新客户跟进失败", h.log)
		return
	}

	response.Success(c, result)
}

// UpdateStatus 修改客户跟进状态。
func (h *Handler) UpdateStatus(c *gin.Context) {
	actor, ok := currentActor(c, h.log)
	if !ok {
		return
	}

	followUpID, err := ParseFollowUpID(c.Param("id"))
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	if err := h.service.UpdateStatus(actor, followUpID, req.Status); err != nil {
		writeError(c, err, "更新客户跟进状态失败", h.log)
		return
	}

	response.Success(c, gin.H{
		"id":     followUpID,
		"status": req.Status,
	})
}

func currentActor(c *gin.Context, log *zap.Logger) (datascope.Actor, bool) {
	actor, ok := middleware.CurrentActor(c)
	if !ok {
		response.Error(c, apperror.Unauthorized("请先登录"), log)
		return datascope.Actor{}, false
	}
	return actor, true
}

func writeError(c *gin.Context, err error, fallbackMessage string, log *zap.Logger) {
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		response.Error(c, appErr, log)
		return
	}

	response.Error(c, apperror.Internal(fallbackMessage, err), log)
}
