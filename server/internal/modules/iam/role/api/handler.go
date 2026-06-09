// Package api 提供角色模块的 HTTP 请求处理器与路由定义。
package api

import (
	roleapp "ez-admin-gin/server/internal/modules/iam/role/application"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	httpx "ez-admin-gin/server/internal/pkg/httpx"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	service *roleapp.Service
	log     *zap.Logger
}

func NewHandler(service *roleapp.Service, log *zap.Logger) *Handler {
	return &Handler{service: service, log: log}
}

// List godoc
// @Summary      查询角色列表
// @Tags         IAM / 角色管理
// @Accept       json
// @Produce      json
// @Param        keyword    query     string  false  "关键词"
// @Param        status     query     int     false  "状态"
// @Success      200  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/roles [get]
func (h *Handler) List(c *gin.Context) {
	var query ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httpx.Error(c, errorsx.BadRequest("查询参数不正确"), h.log)
		return
	}

	result, err := h.service.List(query)
	if err != nil {
		httpx.WriteError(c, err, "查询角色列表失败", h.log)
		return
	}

	httpx.Success(c, result)
}

// Create godoc
// @Summary      创建角色
// @Tags         IAM / 角色管理
// @Accept       json
// @Produce      json
// @Param        body  body  CreateRequest  true  "角色信息"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/roles [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.Create(req)
	if err != nil {
		httpx.WriteError(c, err, "创建角色失败", h.log)
		return
	}

	httpx.Success(c, result)
}

// Update godoc
// @Summary      更新角色
// @Tags         IAM / 角色管理
// @Accept       json
// @Produce      json
// @Param        id    path      uint            true  "角色 ID"
// @Param        body  body      UpdateRequest   true  "角色信息"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/roles/{id}/update [post]
func (h *Handler) Update(c *gin.Context) {
	roleID, ok := httpx.UintIDParam(c, "id", "角色 ID", h.log)
	if !ok {
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.Update(roleID, req)
	if err != nil {
		httpx.WriteError(c, err, "更新角色失败", h.log)
		return
	}

	httpx.Success(c, result)
}

// UpdateStatus godoc
// @Summary      更新角色状态
// @Tags         IAM / 角色管理
// @Accept       json
// @Produce      json
// @Param        id    path      uint                  true  "角色 ID"
// @Param        body  body      UpdateStatusRequest   true  "状态"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/roles/{id}/status [post]
func (h *Handler) UpdateStatus(c *gin.Context) {
	roleID, ok := httpx.UintIDParam(c, "id", "角色 ID", h.log)
	if !ok {
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	if err := h.service.UpdateStatus(roleID, req.Status); err != nil {
		httpx.WriteError(c, err, "更新角色状态失败", h.log)
		return
	}

	httpx.Success(c, gin.H{"id": roleID, "status": req.Status})
}

// UpdatePermissions godoc
// @Summary      更新角色接口权限
// @Tags         IAM / 角色管理
// @Accept       json
// @Produce      json
// @Param        id    path      uint                      true  "角色 ID"
// @Param        body  body      UpdatePermissionsRequest  true  "权限列表"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/roles/{id}/permissions [post]
func (h *Handler) UpdatePermissions(c *gin.Context) {
	roleID, ok := httpx.UintIDParam(c, "id", "角色 ID", h.log)
	if !ok {
		return
	}

	var req UpdatePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	apiIDs, roleCode, err := h.service.UpdatePermissions(roleID, req.APIIDs)
	if err != nil {
		httpx.WriteError(c, err, "更新角色接口权限失败", h.log)
		return
	}

	httpx.Success(c, gin.H{"id": roleID, "code": roleCode, "api_ids": apiIDs})
}

// UpdateMenus godoc
// @Summary      更新角色菜单权限
// @Tags         IAM / 角色管理
// @Accept       json
// @Produce      json
// @Param        id    path      uint                  true  "角色 ID"
// @Param        body  body      UpdateMenusRequest    true  "菜单 ID 列表"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/roles/{id}/menus [post]
func (h *Handler) UpdateMenus(c *gin.Context) {
	roleID, ok := httpx.UintIDParam(c, "id", "角色 ID", h.log)
	if !ok {
		return
	}

	var req UpdateMenusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	menuIDs, err := h.service.UpdateMenus(roleID, req.MenuIDs)
	if err != nil {
		httpx.WriteError(c, err, "更新角色菜单权限失败", h.log)
		return
	}

	httpx.Success(c, gin.H{"id": roleID, "menu_ids": menuIDs})
}

// Delete godoc
// @Summary      删除角色
// @Tags         IAM / 角色管理
// @Accept       json
// @Produce      json
// @Param        id    path      uint  true  "角色 ID"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/roles/{id}/delete [post]
func (h *Handler) Delete(c *gin.Context) {
	roleID, ok := httpx.UintIDParam(c, "id", "角色 ID", h.log)
	if !ok {
		return
	}

	if err := h.service.Delete(roleID); err != nil {
		httpx.WriteError(c, err, "删除角色失败", h.log)
		return
	}

	httpx.Success(c, gin.H{"id": roleID})
}
