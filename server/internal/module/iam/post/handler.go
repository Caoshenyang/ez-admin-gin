package post

import (
	"errors"
	"strconv"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	service *Service
	log     *zap.Logger
}

func NewHandler(service *Service, log *zap.Logger) *Handler {
	return &Handler{service: service, log: log}
}

func (h *Handler) List(c *gin.Context) {
	var query ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, apperror.BadRequest("查询参数不正确"), h.log)
		return
	}

	result, err := h.service.List(query)
	if err != nil {
		writeError(c, err, "查询岗位列表失败", h.log)
		return
	}

	response.Success(c, result)
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.Create(req)
	if err != nil {
		writeError(c, err, "创建岗位失败", h.log)
		return
	}

	response.Success(c, result)
}

func (h *Handler) Update(c *gin.Context) {
	postID, ok := postIDParam(c, h.log)
	if !ok {
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.Update(postID, req)
	if err != nil {
		writeError(c, err, "更新岗位失败", h.log)
		return
	}

	response.Success(c, result)
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	postID, ok := postIDParam(c, h.log)
	if !ok {
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	if err := h.service.UpdateStatus(postID, req.Status); err != nil {
		writeError(c, err, "更新岗位状态失败", h.log)
		return
	}

	response.Success(c, gin.H{"id": postID, "status": req.Status})
}

func postIDParam(c *gin.Context, log *zap.Logger) (uint, bool) {
	rawID := c.Param("id")
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 {
		response.Error(c, apperror.BadRequest("岗位 ID 不正确"), log)
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
