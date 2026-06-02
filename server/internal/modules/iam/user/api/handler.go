// Package api 提供用户模块的 HTTP 请求处理器与路由定义。
package api

import (
	userapp "ez-admin-gin/server/internal/modules/iam/user/application"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	httpx "ez-admin-gin/server/internal/pkg/httpx"
	"ez-admin-gin/server/internal/platform/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	service *userapp.Service
	log     *zap.Logger
}

func NewHandler(service *userapp.Service, log *zap.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}

// List godoc
// @Summary      查询用户列表
// @Tags         IAM / 用户管理
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "页码"
// @Param        page_size  query     int     false  "每页条数"
// @Param        keyword    query     string  false  "关键词"
// @Param        status     query     int     false  "状态"
// @Param        department_id  query  int    false  "部门 ID"
// @Param        role_id    query     int     false  "角色 ID"
// @Success      200  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/users [get]
func (h *Handler) List(c *gin.Context) {
	actor, ok := httpx.CurrentActor(c, h.log)
	if !ok {
		return
	}

	var query ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httpx.Error(c, errorsx.BadRequest("查询参数不正确"), h.log)
		return
	}

	result, err := h.service.List(actor, query)
	if err != nil {
		httpx.WriteError(c, err, "查询用户列表失败", h.log)
		return
	}

	httpx.Success(c, result)
}

// Create godoc
// @Summary      创建用户
// @Tags         IAM / 用户管理
// @Accept       json
// @Produce      json
// @Param        body  body  CreateRequest  true  "用户信息"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/users [post]
func (h *Handler) Create(c *gin.Context) {
	actor, ok := httpx.CurrentActor(c, h.log)
	if !ok {
		return
	}

	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.Create(actor, req)
	if err != nil {
		httpx.WriteError(c, err, "创建用户失败", h.log)
		return
	}

	httpx.Success(c, result)
}

// Update godoc
// @Summary      更新用户
// @Tags         IAM / 用户管理
// @Accept       json
// @Produce      json
// @Param        id    path      uint            true  "用户 ID"
// @Param        body  body      UpdateRequest   true  "用户信息"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/users/{id}/update [post]
func (h *Handler) Update(c *gin.Context) {
	actor, ok := httpx.CurrentActor(c, h.log)
	if !ok {
		return
	}

	userID, ok := httpx.UintIDParam(c, "id", "用户 ID", h.log)
	if !ok {
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	currentUserID, _ := middleware.CurrentUserID(c)
	result, err := h.service.Update(actor, userID, currentUserID, req)
	if err != nil {
		httpx.WriteError(c, err, "更新用户失败", h.log)
		return
	}

	httpx.Success(c, result)
}

// UpdateStatus godoc
// @Summary      更新用户状态
// @Tags         IAM / 用户管理
// @Accept       json
// @Produce      json
// @Param        id    path      uint                  true  "用户 ID"
// @Param        body  body      UpdateStatusRequest   true  "状态"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/users/{id}/status [post]
func (h *Handler) UpdateStatus(c *gin.Context) {
	actor, ok := httpx.CurrentActor(c, h.log)
	if !ok {
		return
	}

	userID, ok := httpx.UintIDParam(c, "id", "用户 ID", h.log)
	if !ok {
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	currentUserID, _ := middleware.CurrentUserID(c)
	if err := h.service.UpdateStatus(actor, userID, currentUserID, uint(req.Status)); err != nil {
		httpx.WriteError(c, err, "更新用户状态失败", h.log)
		return
	}

	httpx.Success(c, gin.H{
		"id":     userID,
		"status": req.Status,
	})
}

// UpdateRoles godoc
// @Summary      更新用户角色
// @Tags         IAM / 用户管理
// @Accept       json
// @Produce      json
// @Param        id    path      uint                  true  "用户 ID"
// @Param        body  body      UpdateRolesRequest    true  "角色 ID 列表"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/users/{id}/roles [post]
func (h *Handler) UpdateRoles(c *gin.Context) {
	actor, ok := httpx.CurrentActor(c, h.log)
	if !ok {
		return
	}

	userID, ok := httpx.UintIDParam(c, "id", "用户 ID", h.log)
	if !ok {
		return
	}

	var req UpdateRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	currentUserID, _ := middleware.CurrentUserID(c)
	roleIDs, err := h.service.UpdateRoles(actor, userID, currentUserID, req.RoleIDs)
	if err != nil {
		httpx.WriteError(c, err, "更新用户角色失败", h.log)
		return
	}

	httpx.Success(c, gin.H{
		"id":       userID,
		"role_ids": roleIDs,
	})
}
