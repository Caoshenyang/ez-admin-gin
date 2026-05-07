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

	permissions, roleCode, err := h.service.UpdatePermissions(roleID, req.Permissions)
	if err != nil {
		httpx.WriteError(c, err, "更新角色接口权限失败", h.log)
		return
	}

	httpx.Success(c, gin.H{"id": roleID, "code": roleCode, "permissions": permissions})
}

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
