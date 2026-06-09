// Package api 提供部门模块的 HTTP 请求处理器与路由定义。
package api

import (
	departmentapp "ez-admin-gin/server/internal/modules/iam/department/application"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	httpx "ez-admin-gin/server/internal/pkg/httpx"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	service *departmentapp.Service
	log     *zap.Logger
}

func NewHandler(service *departmentapp.Service, log *zap.Logger) *Handler {
	return &Handler{service: service, log: log}
}

// List godoc
// @Summary      查询部门列表
// @Tags         IAM / 部门管理
// @Accept       json
// @Produce      json
// @Param        keyword    query     string  false  "关键词"
// @Param        status     query     int     false  "状态"
// @Success      200  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/departments [get]
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
		httpx.WriteError(c, err, "查询部门列表失败", h.log)
		return
	}

	httpx.Success(c, result)
}

// Create godoc
// @Summary      创建部门
// @Tags         IAM / 部门管理
// @Accept       json
// @Produce      json
// @Param        body  body  CreateRequest  true  "部门信息"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/departments [post]
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
		httpx.WriteError(c, err, "创建部门失败", h.log)
		return
	}

	httpx.Success(c, result)
}

// Update godoc
// @Summary      更新部门
// @Tags         IAM / 部门管理
// @Accept       json
// @Produce      json
// @Param        id    path      uint            true  "部门 ID"
// @Param        body  body      UpdateRequest   true  "部门信息"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/departments/{id}/update [post]
func (h *Handler) Update(c *gin.Context) {
	actor, ok := httpx.CurrentActor(c, h.log)
	if !ok {
		return
	}

	departmentID, ok := httpx.UintIDParam(c, "id", "部门 ID", h.log)
	if !ok {
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.Update(actor, departmentID, req)
	if err != nil {
		httpx.WriteError(c, err, "更新部门失败", h.log)
		return
	}

	httpx.Success(c, result)
}

// UpdateStatus godoc
// @Summary      更新部门状态
// @Tags         IAM / 部门管理
// @Accept       json
// @Produce      json
// @Param        id    path      uint                  true  "部门 ID"
// @Param        body  body      UpdateStatusRequest   true  "状态"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/departments/{id}/status [post]
func (h *Handler) UpdateStatus(c *gin.Context) {
	actor, ok := httpx.CurrentActor(c, h.log)
	if !ok {
		return
	}

	departmentID, ok := httpx.UintIDParam(c, "id", "部门 ID", h.log)
	if !ok {
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	if err := h.service.UpdateStatus(actor, departmentID, req.Status); err != nil {
		httpx.WriteError(c, err, "更新部门状态失败", h.log)
		return
	}

	httpx.Success(c, gin.H{"id": departmentID, "status": req.Status})
}

// Delete godoc
// @Summary      删除部门
// @Tags         IAM / 部门管理
// @Accept       json
// @Produce      json
// @Param        id    path      uint  true  "部门 ID"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/departments/{id}/delete [post]
func (h *Handler) Delete(c *gin.Context) {
	actor, ok := httpx.CurrentActor(c, h.log)
	if !ok {
		return
	}

	departmentID, ok := httpx.UintIDParam(c, "id", "部门 ID", h.log)
	if !ok {
		return
	}

	if err := h.service.Delete(actor, departmentID); err != nil {
		httpx.WriteError(c, err, "删除部门失败", h.log)
		return
	}

	httpx.Success(c, gin.H{"id": departmentID})
}
