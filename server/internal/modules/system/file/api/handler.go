package api

import (
	fileapp "ez-admin-gin/server/internal/modules/system/file/application"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	httpx "ez-admin-gin/server/internal/pkg/httpx"
	"ez-admin-gin/server/internal/platform/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	service *fileapp.Service
	log     *zap.Logger
}

func NewHandler(service *fileapp.Service, log *zap.Logger) *Handler {
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
		httpx.WriteError(c, err, "查询文件列表失败", h.log)
		return
	}

	httpx.Success(c, result)
}

func (h *Handler) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		httpx.Error(c, errorsx.BadRequest("请选择要上传的文件"), h.log)
		return
	}

	uploaderID, _ := middleware.CurrentUserID(c)
	result, err := h.service.Upload(c.Request.Context(), uploaderID, fileHeader)
	if err != nil {
		httpx.WriteError(c, err, "上传文件失败", h.log)
		return
	}

	httpx.Success(c, result)
}
