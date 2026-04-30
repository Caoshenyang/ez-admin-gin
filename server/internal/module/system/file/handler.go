package file

import (
	"errors"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/middleware"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 负责文件模块的 HTTP 协议层绑定与输出。
type Handler struct {
	service *Service
	log     *zap.Logger
}

// NewHandler 创建文件 Handler。
func NewHandler(service *Service, log *zap.Logger) *Handler {
	return &Handler{service: service, log: log}
}

// List 返回文件分页列表。
func (h *Handler) List(c *gin.Context) {
	var query ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, apperror.BadRequest("查询参数不正确"), h.log)
		return
	}

	result, err := h.service.List(query)
	if err != nil {
		writeError(c, err, "查询文件列表失败", h.log)
		return
	}

	response.Success(c, result)
}

// Upload 上传文件并写入记录。
func (h *Handler) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Error(c, apperror.BadRequest("请选择要上传的文件"), h.log)
		return
	}

	uploaderID, _ := middleware.CurrentUserID(c)
	result, err := h.service.Upload(c.Request.Context(), uploaderID, fileHeader)
	if err != nil {
		writeError(c, err, "上传文件失败", h.log)
		return
	}

	response.Success(c, result)
}

func writeError(c *gin.Context, err error, fallbackMessage string, log *zap.Logger) {
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		response.Error(c, appErr, log)
		return
	}

	response.Error(c, apperror.Internal(fallbackMessage, err), log)
}
