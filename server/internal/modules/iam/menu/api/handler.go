package api

import (
	menuapp "ez-admin-gin/server/internal/modules/iam/menu/application"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	httpx "ez-admin-gin/server/internal/pkg/httpx"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	service *menuapp.Service
	log     *zap.Logger
}

func NewHandler(service *menuapp.Service, log *zap.Logger) *Handler {
	return &Handler{service: service, log: log}
}

// List godoc
// @Summary      查询菜单树
// @Tags         IAM / 菜单管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /api/v1/system/menus [get]
func (h *Handler) List(c *gin.Context) {
	result, err := h.service.List()
	if err != nil {
		httpx.WriteError(c, err, "查询菜单树失败", h.log)
		return
	}

	httpx.Success(c, result)
}

// Create godoc
// @Summary      创建菜单
// @Tags         IAM / 菜单管理
// @Accept       json
// @Produce      json
// @Param        body  body  CreateRequest  true  "菜单信息"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /api/v1/system/menus [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.Create(req)
	if err != nil {
		httpx.WriteError(c, err, "创建菜单失败", h.log)
		return
	}

	httpx.Success(c, result)
}

// Update godoc
// @Summary      更新菜单
// @Tags         IAM / 菜单管理
// @Accept       json
// @Produce      json
// @Param        id    path      uint            true  "菜单 ID"
// @Param        body  body      UpdateRequest   true  "菜单信息"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /api/v1/system/menus/{id}/update [post]
func (h *Handler) Update(c *gin.Context) {
	menuID, ok := httpx.UintIDParam(c, "id", "菜单 ID", h.log)
	if !ok {
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.Update(menuID, req)
	if err != nil {
		httpx.WriteError(c, err, "更新菜单失败", h.log)
		return
	}

	httpx.Success(c, result)
}

// UpdateStatus godoc
// @Summary      更新菜单状态
// @Tags         IAM / 菜单管理
// @Accept       json
// @Produce      json
// @Param        id    path      uint                  true  "菜单 ID"
// @Param        body  body      UpdateStatusRequest   true  "状态"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /api/v1/system/menus/{id}/status [post]
func (h *Handler) UpdateStatus(c *gin.Context) {
	menuID, ok := httpx.UintIDParam(c, "id", "菜单 ID", h.log)
	if !ok {
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	if err := h.service.UpdateStatus(menuID, req.Status); err != nil {
		httpx.WriteError(c, err, "更新菜单状态失败", h.log)
		return
	}

	httpx.Success(c, gin.H{"id": menuID, "status": req.Status})
}

// Delete godoc
// @Summary      删除菜单
// @Tags         IAM / 菜单管理
// @Accept       json
// @Produce      json
// @Param        id    path      uint  true  "菜单 ID"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /api/v1/system/menus/{id}/delete [post]
func (h *Handler) Delete(c *gin.Context) {
	menuID, ok := httpx.UintIDParam(c, "id", "菜单 ID", h.log)
	if !ok {
		return
	}

	if err := h.service.Delete(menuID); err != nil {
		httpx.WriteError(c, err, "删除菜单失败", h.log)
		return
	}

	httpx.Success(c, gin.H{"id": menuID})
}
