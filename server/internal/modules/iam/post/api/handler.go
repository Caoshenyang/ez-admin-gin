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
