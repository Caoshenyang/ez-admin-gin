package api

import (
	attachmentapp "ez-admin-gin/server/internal/modules/system/attachment/application"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	httpx "ez-admin-gin/server/internal/pkg/httpx"
	"ez-admin-gin/server/internal/platform/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	service *attachmentapp.Service
	log     *zap.Logger
}

func NewHandler(service *attachmentapp.Service, log *zap.Logger) *Handler {
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
		httpx.WriteError(c, err, "查询附件列表失败", h.log)
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

	var req CreateRequest
	if err := c.ShouldBind(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	uploaderID, _ := middleware.CurrentUserID(c)
	result, err := h.service.CreateByUpload(c.Request.Context(), uploaderID, fileHeader, req)
	if err != nil {
		httpx.WriteError(c, err, "创建附件失败", h.log)
		return
	}

	httpx.Success(c, result)
}

func (h *Handler) Update(c *gin.Context) {
	id, ok := httpx.UintIDParam(c, "id", "附件 ID", h.log)
	if !ok {
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.Update(id, req)
	if err != nil {
		httpx.WriteError(c, err, "更新附件失败", h.log)
		return
	}

	httpx.Success(c, result)
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	id, ok := httpx.UintIDParam(c, "id", "附件 ID", h.log)
	if !ok {
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	if err := h.service.UpdateStatus(id, req.Status); err != nil {
		httpx.WriteError(c, err, "更新附件状态失败", h.log)
		return
	}

	httpx.Success(c, gin.H{"id": id, "status": req.Status})
}
