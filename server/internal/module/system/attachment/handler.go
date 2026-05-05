package attachment

import (
	"errors"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/middleware"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 负责附件中心的 HTTP 协议层绑定与输出。
type Handler struct {
	service *Service
	log     *zap.Logger
}

// NewHandler 创建附件中心 Handler。
func NewHandler(service *Service, log *zap.Logger) *Handler {
	return &Handler{service: service, log: log}
}

// List 返回附件中心分页列表。
func (h *Handler) List(c *gin.Context) {
	var query ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, apperror.BadRequest("查询参数不正确"), h.log)
		return
	}

	result, err := h.service.List(query)
	if err != nil {
		writeError(c, err, "查询附件列表失败", h.log)
		return
	}

	response.Success(c, result)
}

// Upload 上传文件并创建附件中心记录。
func (h *Handler) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Error(c, apperror.BadRequest("请选择要上传的文件"), h.log)
		return
	}

	var req CreateRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	uploaderID, _ := middleware.CurrentUserID(c)
	result, err := h.service.CreateByUpload(c.Request.Context(), uploaderID, fileHeader, req)
	if err != nil {
		writeError(c, err, "创建附件失败", h.log)
		return
	}

	response.Success(c, result)
}

// Update 修改附件元数据。
func (h *Handler) Update(c *gin.Context) {
	id, err := ParseAttachmentID(c.Param("id"))
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.Update(id, req)
	if err != nil {
		writeError(c, err, "更新附件失败", h.log)
		return
	}

	response.Success(c, result)
}

// UpdateStatus 单独修改附件状态。
func (h *Handler) UpdateStatus(c *gin.Context) {
	id, err := ParseAttachmentID(c.Param("id"))
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	if err := h.service.UpdateStatus(id, req.Status); err != nil {
		writeError(c, err, "更新附件状态失败", h.log)
		return
	}

	response.Success(c, gin.H{
		"id":     id,
		"status": req.Status,
	})
}

func writeError(c *gin.Context, err error, fallbackMessage string, log *zap.Logger) {
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		response.Error(c, appErr, log)
		return
	}

	response.Error(c, apperror.Internal(fallbackMessage, err), log)
}
