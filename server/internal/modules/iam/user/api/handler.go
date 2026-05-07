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
