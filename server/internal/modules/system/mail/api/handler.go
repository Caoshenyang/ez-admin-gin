// Package api 提供邮件模块的 HTTP 请求处理器与路由定义。
package api

import (
	mailapp "ez-admin-gin/server/internal/modules/system/mail/application"
	maildomain "ez-admin-gin/server/internal/modules/system/mail/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	httpx "ez-admin-gin/server/internal/pkg/httpx"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	service *mailapp.Service
	log     *zap.Logger
}

func NewHandler(service *mailapp.Service, log *zap.Logger) *Handler {
	return &Handler{service: service, log: log}
}

// ListAccounts godoc
// @Summary      查询系统邮箱列表
// @Tags         System / 邮件管理
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "页码"
// @Param        page_size  query     int     false  "每页条数"
// @Param        keyword    query     string  false  "关键词"
// @Param        status     query     int     false  "状态"
// @Success      200  {object}  httpx.Body{data=AccountListResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/mail/accounts [get]
func (h *Handler) ListAccounts(c *gin.Context) {
	var query AccountListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httpx.Error(c, errorsx.BadRequest("查询参数不正确"), h.log)
		return
	}
	result, err := h.service.ListAccounts(query)
	if err != nil {
		httpx.WriteError(c, err, "查询系统邮箱列表失败", h.log)
		return
	}
	httpx.Success(c, result)
}

// CreateAccount godoc
// @Summary      创建系统邮箱
// @Tags         System / 邮件管理
// @Accept       json
// @Produce      json
// @Param        body  body  CreateAccountRequest  true  "邮箱参数"
// @Success      200  {object}  httpx.Body{data=AccountResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/mail/accounts [post]
func (h *Handler) CreateAccount(c *gin.Context) {
	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}
	result, err := h.service.CreateAccount(req)
	if err != nil {
		httpx.WriteError(c, err, "创建系统邮箱失败", h.log)
		return
	}
	httpx.Success(c, result)
}

// UpdateAccount godoc
// @Summary      更新系统邮箱
// @Tags         System / 邮件管理
// @Accept       json
// @Produce      json
// @Param        id    path  uint                  true  "邮箱 ID"
// @Param        body  body  UpdateAccountRequest  true  "邮箱参数"
// @Success      200  {object}  httpx.Body{data=AccountResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/mail/accounts/{id}/update [post]
func (h *Handler) UpdateAccount(c *gin.Context) {
	accountID, ok := httpx.UintIDParam(c, "id", "邮箱 ID", h.log)
	if !ok {
		return
	}
	var req UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}
	result, err := h.service.UpdateAccount(accountID, req)
	if err != nil {
		httpx.WriteError(c, err, "更新系统邮箱失败", h.log)
		return
	}
	httpx.Success(c, result)
}

// UpdateAccountStatus godoc
// @Summary      更新系统邮箱状态
// @Tags         System / 邮件管理
// @Accept       json
// @Produce      json
// @Param        id    path  uint                        true  "邮箱 ID"
// @Param        body  body  UpdateAccountStatusRequest  true  "状态参数"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/mail/accounts/{id}/status [post]
func (h *Handler) UpdateAccountStatus(c *gin.Context) {
	accountID, ok := httpx.UintIDParam(c, "id", "邮箱 ID", h.log)
	if !ok {
		return
	}
	var req UpdateAccountStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}
	if err := h.service.UpdateAccountStatus(accountID, req.Status); err != nil {
		httpx.WriteError(c, err, "更新系统邮箱状态失败", h.log)
		return
	}
	httpx.Success(c, gin.H{"id": accountID, "status": req.Status})
}

// DeleteAccount godoc
// @Summary      删除系统邮箱
// @Tags         System / 邮件管理
// @Accept       json
// @Produce      json
// @Param        id  path  uint  true  "邮箱 ID"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/mail/accounts/{id}/delete [post]
func (h *Handler) DeleteAccount(c *gin.Context) {
	accountID, ok := httpx.UintIDParam(c, "id", "邮箱 ID", h.log)
	if !ok {
		return
	}
	if err := h.service.DeleteAccount(accountID); err != nil {
		httpx.WriteError(c, err, "删除系统邮箱失败", h.log)
		return
	}
	httpx.Success(c, gin.H{"id": accountID})
}

// TestAccount godoc
// @Summary      测试系统邮箱
// @Tags         System / 邮件管理
// @Accept       json
// @Produce      json
// @Param        id    path  uint                true  "邮箱 ID"
// @Param        body  body  TestAccountRequest  true  "测试参数"
// @Success      200  {object}  httpx.Body{data=SendMailResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/mail/accounts/{id}/test [post]
func (h *Handler) TestAccount(c *gin.Context) {
	accountID, ok := httpx.UintIDParam(c, "id", "邮箱 ID", h.log)
	if !ok {
		return
	}
	var req TestAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}
	result, err := h.service.TestAccount(accountID, req)
	if err != nil {
		httpx.WriteError(c, err, "测试系统邮箱失败", h.log)
		return
	}
	httpx.Success(c, result)
}

// ListTemplates godoc
// @Summary      查询邮件模板列表
// @Tags         System / 邮件管理
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "页码"
// @Param        page_size  query     int     false  "每页条数"
// @Param        keyword    query     string  false  "关键词"
// @Param        status     query     int     false  "状态"
// @Success      200  {object}  httpx.Body{data=TemplateListResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/mail/templates [get]
func (h *Handler) ListTemplates(c *gin.Context) {
	var query TemplateListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httpx.Error(c, errorsx.BadRequest("查询参数不正确"), h.log)
		return
	}
	result, err := h.service.ListTemplates(query)
	if err != nil {
		httpx.WriteError(c, err, "查询邮件模板列表失败", h.log)
		return
	}
	httpx.Success(c, result)
}

// CreateTemplate godoc
// @Summary      创建邮件模板
// @Tags         System / 邮件管理
// @Accept       json
// @Produce      json
// @Param        body  body  CreateTemplateRequest  true  "模板参数"
// @Success      200  {object}  httpx.Body{data=TemplateResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/mail/templates [post]
func (h *Handler) CreateTemplate(c *gin.Context) {
	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}
	result, err := h.service.CreateTemplate(req)
	if err != nil {
		httpx.WriteError(c, err, "创建邮件模板失败", h.log)
		return
	}
	httpx.Success(c, result)
}

// UpdateTemplate godoc
// @Summary      更新邮件模板
// @Tags         System / 邮件管理
// @Accept       json
// @Produce      json
// @Param        id    path  uint                   true  "模板 ID"
// @Param        body  body  UpdateTemplateRequest  true  "模板参数"
// @Success      200  {object}  httpx.Body{data=TemplateResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/mail/templates/{id}/update [post]
func (h *Handler) UpdateTemplate(c *gin.Context) {
	templateID, ok := httpx.UintIDParam(c, "id", "模板 ID", h.log)
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
		httpx.WriteError(c, err, "更新邮件模板失败", h.log)
		return
	}
	httpx.Success(c, result)
}

// UpdateTemplateStatus godoc
// @Summary      更新邮件模板状态
// @Tags         System / 邮件管理
// @Accept       json
// @Produce      json
// @Param        id    path  uint                         true  "模板 ID"
// @Param        body  body  UpdateTemplateStatusRequest  true  "状态参数"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/mail/templates/{id}/status [post]
func (h *Handler) UpdateTemplateStatus(c *gin.Context) {
	templateID, ok := httpx.UintIDParam(c, "id", "模板 ID", h.log)
	if !ok {
		return
	}
	var req UpdateTemplateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}
	if err := h.service.UpdateTemplateStatus(templateID, req.Status); err != nil {
		httpx.WriteError(c, err, "更新邮件模板状态失败", h.log)
		return
	}
	httpx.Success(c, gin.H{"id": templateID, "status": req.Status})
}

// DeleteTemplate godoc
// @Summary      删除邮件模板
// @Tags         System / 邮件管理
// @Accept       json
// @Produce      json
// @Param        id  path  uint  true  "模板 ID"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/mail/templates/{id}/delete [post]
func (h *Handler) DeleteTemplate(c *gin.Context) {
	templateID, ok := httpx.UintIDParam(c, "id", "模板 ID", h.log)
	if !ok {
		return
	}
	if err := h.service.DeleteTemplate(templateID); err != nil {
		httpx.WriteError(c, err, "删除邮件模板失败", h.log)
		return
	}
	httpx.Success(c, gin.H{"id": templateID})
}

// RenderTemplate godoc
// @Summary      渲染邮件模板
// @Tags         System / 邮件管理
// @Accept       json
// @Produce      json
// @Param        id    path  uint                   true  "模板 ID"
// @Param        body  body  RenderTemplateRequest  true  "变量参数"
// @Success      200  {object}  httpx.Body{data=RenderTemplateResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/mail/templates/{id}/render [post]
func (h *Handler) RenderTemplate(c *gin.Context) {
	templateID, ok := httpx.UintIDParam(c, "id", "模板 ID", h.log)
	if !ok {
		return
	}
	var req RenderTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}
	result, err := h.service.RenderTemplate(templateID, req)
	if err != nil {
		httpx.WriteError(c, err, "渲染邮件模板失败", h.log)
		return
	}
	httpx.Success(c, result)
}

// Send godoc
// @Summary      发送邮件
// @Tags         System / 邮件管理
// @Accept       json
// @Produce      json
// @Param        body  body  SendMailRequest  true  "邮件参数"
// @Success      200  {object}  httpx.Body{data=SendMailResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/mail/send [post]
func (h *Handler) Send(c *gin.Context) {
	var req SendMailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}
	result, err := h.service.Send(req)
	if err != nil {
		httpx.WriteError(c, err, "发送邮件失败", h.log)
		return
	}
	httpx.Success(c, result)
}

// ListLogs godoc
// @Summary      查询邮件发送日志
// @Tags         System / 邮件管理
// @Accept       json
// @Produce      json
// @Param        page           query     int     false  "页码"
// @Param        page_size      query     int     false  "每页条数"
// @Param        keyword        query     string  false  "关键词"
// @Param        status         query     int     false  "发送状态"
// @Param        account_id     query     int     false  "邮箱 ID"
// @Param        template_code  query     string  false  "模板编码"
// @Success      200  {object}  httpx.Body{data=LogListResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/mail/logs [get]
func (h *Handler) ListLogs(c *gin.Context) {
	var query LogListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httpx.Error(c, errorsx.BadRequest("查询参数不正确"), h.log)
		return
	}
	result, err := h.service.ListLogs(query)
	if err != nil {
		httpx.WriteError(c, err, "查询邮件发送日志失败", h.log)
		return
	}
	httpx.Success(c, result)
}

var _ = maildomain.PermissionAccountList
