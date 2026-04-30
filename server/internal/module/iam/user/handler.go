package user

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

// Handler 负责把 HTTP 协议层请求转成用户服务调用。
type Handler struct {
	service *Service
	log     *zap.Logger
}

// NewHandler 创建用户模块 Handler。
func NewHandler(service *Service, log *zap.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}

// List 返回当前数据范围内的用户分页列表。
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
		writeError(c, err, "查询用户列表失败", h.log)
		return
	}

	response.Success(c, result)
}

// Create 创建后台用户。
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
		writeError(c, err, "创建用户失败", h.log)
		return
	}

	response.Success(c, result)
}

// Update 编辑用户基础信息。
func (h *Handler) Update(c *gin.Context) {
	actor, ok := currentActor(c, h.log)
	if !ok {
		return
	}

	userID, ok := userIDParam(c, h.log)
	if !ok {
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	currentUserID, _ := middleware.CurrentUserID(c)
	result, err := h.service.Update(actor, userID, currentUserID, req)
	if err != nil {
		writeError(c, err, "更新用户失败", h.log)
		return
	}

	response.Success(c, result)
}

// UpdateStatus 修改用户启用状态。
func (h *Handler) UpdateStatus(c *gin.Context) {
	actor, ok := currentActor(c, h.log)
	if !ok {
		return
	}

	userID, ok := userIDParam(c, h.log)
	if !ok {
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	currentUserID, _ := middleware.CurrentUserID(c)
	if err := h.service.UpdateStatus(actor, userID, currentUserID, uint(req.Status)); err != nil {
		writeError(c, err, "更新用户状态失败", h.log)
		return
	}

	response.Success(c, gin.H{
		"id":     userID,
		"status": req.Status,
	})
}

// UpdateRoles 更新用户绑定的角色。
func (h *Handler) UpdateRoles(c *gin.Context) {
	actor, ok := currentActor(c, h.log)
	if !ok {
		return
	}

	userID, ok := userIDParam(c, h.log)
	if !ok {
		return
	}

	var req UpdateRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	currentUserID, _ := middleware.CurrentUserID(c)
	roleIDs, err := h.service.UpdateRoles(actor, userID, currentUserID, req.RoleIDs)
	if err != nil {
		writeError(c, err, "更新用户角色失败", h.log)
		return
	}

	response.Success(c, gin.H{
		"id":       userID,
		"role_ids": roleIDs,
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

func userIDParam(c *gin.Context, log *zap.Logger) (uint, bool) {
	rawID := c.Param("id")
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 {
		response.Error(c, apperror.BadRequest("用户 ID 不正确"), log)
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
