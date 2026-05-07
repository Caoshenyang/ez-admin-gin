package api

import (
	postapp "ez-admin-gin/server/internal/modules/iam/post/application"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	httpx "ez-admin-gin/server/internal/pkg/httpx"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	service *postapp.Service
	log     *zap.Logger
}

func NewHandler(service *postapp.Service, log *zap.Logger) *Handler {
	return &Handler{service: service, log: log}
}

// List godoc
// @Summary      查询岗位列表
// @Tags         IAM / 岗位管理
// @Accept       json
// @Produce      json
// @Param        keyword    query     string  false  "关键词"
// @Param        status     query     int     false  "状态"
// @Success      200  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /api/v1/system/posts [get]
func (h *Handler) List(c *gin.Context) {
	var query ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httpx.Error(c, errorsx.BadRequest("查询参数不正确"), h.log)
		return
	}

	result, err := h.service.List(query)
	if err != nil {
		httpx.WriteError(c, err, "查询岗位列表失败", h.log)
		return
	}

	httpx.Success(c, result)
}

// Create godoc
// @Summary      创建岗位
// @Tags         IAM / 岗位管理
// @Accept       json
// @Produce      json
// @Param        body  body  CreateRequest  true  "岗位信息"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /api/v1/system/posts [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.Create(req)
	if err != nil {
		httpx.WriteError(c, err, "创建岗位失败", h.log)
		return
	}

	httpx.Success(c, result)
}

// Update godoc
// @Summary      更新岗位
// @Tags         IAM / 岗位管理
// @Accept       json
// @Produce      json
// @Param        id    path      uint            true  "岗位 ID"
// @Param        body  body      UpdateRequest   true  "岗位信息"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /api/v1/system/posts/{id}/update [post]
func (h *Handler) Update(c *gin.Context) {
	postID, ok := httpx.UintIDParam(c, "id", "岗位 ID", h.log)
	if !ok {
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.Update(postID, req)
	if err != nil {
		httpx.WriteError(c, err, "更新岗位失败", h.log)
		return
	}

	httpx.Success(c, result)
}

// UpdateStatus godoc
// @Summary      更新岗位状态
// @Tags         IAM / 岗位管理
// @Accept       json
// @Produce      json
// @Param        id    path      uint                  true  "岗位 ID"
// @Param        body  body      UpdateStatusRequest   true  "状态"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /api/v1/system/posts/{id}/status [post]
func (h *Handler) UpdateStatus(c *gin.Context) {
	postID, ok := httpx.UintIDParam(c, "id", "岗位 ID", h.log)
	if !ok {
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	if err := h.service.UpdateStatus(postID, req.Status); err != nil {
		httpx.WriteError(c, err, "更新岗位状态失败", h.log)
		return
	}

	httpx.Success(c, gin.H{"id": postID, "status": req.Status})
}
