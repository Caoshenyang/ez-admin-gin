// Package api 提供消息提醒配置模块的 HTTP 请求处理器与路由定义。
package api

import (
	messageapp "ez-admin-gin/server/internal/modules/system/message/application"
	messagedomain "ez-admin-gin/server/internal/modules/system/message/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	httpx "ez-admin-gin/server/internal/pkg/httpx"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	service *messageapp.Service
	log     *zap.Logger
}

func NewHandler(service *messageapp.Service, log *zap.Logger) *Handler {
	return &Handler{service: service, log: log}
}

// ListTemplates godoc
// @Summary      查询消息模板列表
// @Tags         System / 消息提醒配置
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "页码"
// @Param        page_size  query     int     false  "每页条数"
// @Param        keyword    query     string  false  "关键词"
// @Param        type       query     int     false  "模板类型"
// @Param        status     query     int     false  "状态"
// @Success      200  {object}  httpx.Body{data=TemplateListResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/message-templates [get]
func (h *Handler) ListTemplates(c *gin.Context) {
	var query TemplateListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httpx.Error(c, errorsx.BadRequest("查询参数不正确"), h.log)
		return
	}

	result, err := h.service.ListTemplates(query)
	if err != nil {
		httpx.WriteError(c, err, "查询消息模板列表失败", h.log)
		return
	}

	httpx.Success(c, result)
}

// CreateTemplate godoc
// @Summary      创建消息模板
// @Tags         System / 消息提醒配置
// @Accept       json
// @Produce      json
// @Param        body  body  CreateTemplateRequest  true  "消息模板参数"
// @Success      200  {object}  httpx.Body{data=TemplateResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/message-templates [post]
func (h *Handler) CreateTemplate(c *gin.Context) {
	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.CreateTemplate(req)
	if err != nil {
		httpx.WriteError(c, err, "创建消息模板失败", h.log)
		return
	}

	httpx.Success(c, result)
}

// UpdateTemplate godoc
// @Summary      更新消息模板
// @Tags         System / 消息提醒配置
// @Accept       json
// @Produce      json
// @Param        id    path  uint                   true  "消息模板 ID"
// @Param        body  body  UpdateTemplateRequest  true  "消息模板参数"
// @Success      200  {object}  httpx.Body{data=TemplateResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/message-templates/{id}/update [post]
func (h *Handler) UpdateTemplate(c *gin.Context) {
	templateID, ok := httpx.UintIDParam(c, "id", "消息模板 ID", h.log)
	if !ok {
		return
	}

	var req UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.UpdateTemplate(templateID, req)
	if err != nil {
		httpx.WriteError(c, err, "更新消息模板失败", h.log)
		return
	}

	httpx.Success(c, result)
}

// UpdateTemplateStatus godoc
// @Summary      更新消息模板状态
// @Tags         System / 消息提醒配置
// @Accept       json
// @Produce      json
// @Param        id    path  uint                         true  "消息模板 ID"
// @Param        body  body  UpdateTemplateStatusRequest  true  "状态参数"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/message-templates/{id}/status [post]
func (h *Handler) UpdateTemplateStatus(c *gin.Context) {
	templateID, ok := httpx.UintIDParam(c, "id", "消息模板 ID", h.log)
	if !ok {
		return
	}

	var req UpdateTemplateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	if err := h.service.UpdateTemplateStatus(templateID, req.Status); err != nil {
		httpx.WriteError(c, err, "更新消息模板状态失败", h.log)
		return
	}

	httpx.Success(c, gin.H{"id": templateID, "status": req.Status})
}

// ListReminders godoc
// @Summary      查询提醒规则列表
// @Tags         System / 消息提醒配置
// @Accept       json
// @Produce      json
// @Param        page           query     int     false  "页码"
// @Param        page_size      query     int     false  "每页条数"
// @Param        keyword        query     string  false  "关键词"
// @Param        trigger_event  query     string  false  "触发事件"
// @Param        template_id    query     int     false  "模板 ID"
// @Param        receiver_type  query     int     false  "接收人类型"
// @Param        status         query     int     false  "状态"
// @Success      200  {object}  httpx.Body{data=ReminderListResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/message-reminders [get]
func (h *Handler) ListReminders(c *gin.Context) {
	var query ReminderListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httpx.Error(c, errorsx.BadRequest("查询参数不正确"), h.log)
		return
	}

	result, err := h.service.ListReminders(query)
	if err != nil {
		httpx.WriteError(c, err, "查询提醒规则列表失败", h.log)
		return
	}

	httpx.Success(c, result)
}

// CreateReminder godoc
// @Summary      创建提醒规则
// @Tags         System / 消息提醒配置
// @Accept       json
// @Produce      json
// @Param        body  body  CreateReminderRequest  true  "提醒规则参数"
// @Success      200  {object}  httpx.Body{data=ReminderResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/message-reminders [post]
func (h *Handler) CreateReminder(c *gin.Context) {
	var req CreateReminderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.CreateReminder(req)
	if err != nil {
		httpx.WriteError(c, err, "创建提醒规则失败", h.log)
		return
	}

	httpx.Success(c, result)
}

// UpdateReminder godoc
// @Summary      更新提醒规则
// @Tags         System / 消息提醒配置
// @Accept       json
// @Produce      json
// @Param        id    path  uint                   true  "提醒规则 ID"
// @Param        body  body  UpdateReminderRequest  true  "提醒规则参数"
// @Success      200  {object}  httpx.Body{data=ReminderResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/message-reminders/{id}/update [post]
func (h *Handler) UpdateReminder(c *gin.Context) {
	reminderID, ok := httpx.UintIDParam(c, "id", "提醒规则 ID", h.log)
	if !ok {
		return
	}

	var req UpdateReminderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.UpdateReminder(reminderID, req)
	if err != nil {
		httpx.WriteError(c, err, "更新提醒规则失败", h.log)
		return
	}

	httpx.Success(c, result)
}

// UpdateReminderStatus godoc
// @Summary      更新提醒规则状态
// @Tags         System / 消息提醒配置
// @Accept       json
// @Produce      json
// @Param        id    path  uint                         true  "提醒规则 ID"
// @Param        body  body  UpdateReminderStatusRequest  true  "状态参数"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/message-reminders/{id}/status [post]
func (h *Handler) UpdateReminderStatus(c *gin.Context) {
	reminderID, ok := httpx.UintIDParam(c, "id", "提醒规则 ID", h.log)
	if !ok {
		return
	}

	var req UpdateReminderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	if err := h.service.UpdateReminderStatus(reminderID, req.Status); err != nil {
		httpx.WriteError(c, err, "更新提醒规则状态失败", h.log)
		return
	}

	httpx.Success(c, gin.H{"id": reminderID, "status": req.Status})
}

var _ = messagedomain.PermissionTemplateList
