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

// List godoc
// @Summary      查询附件列表
// @Tags         System / 附件管理
// @Accept       json
// @Produce      json
// @Param        page        query     int     false  "页码"
// @Param        page_size   query     int     false  "每页条数"
// @Param        keyword     query     string  false  "关键词"
// @Param        category    query     string  false  "附件分类"
// @Param        biz_type    query     string  false  "业务类型"
// @Param        ext         query     string  false  "文件后缀"
// @Param        status      query     int     false  "状态"
// @Success      200  {object}  httpx.Body{data=ListResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /api/v1/system/attachments [get]
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

// Upload godoc
// @Summary      上传附件
// @Tags         System / 附件管理
// @Accept       multipart/form-data
// @Produce      json
// @Param        file          formData  file    true   "附件文件"
// @Param        display_name  formData  string  false  "附件名称"
// @Param        category      formData  string  false  "附件分类"
// @Param        biz_type      formData  string  false  "业务类型"
// @Param        status        formData  int     false  "状态"
// @Param        remark        formData  string  false  "备注"
// @Success      200  {object}  httpx.Body{data=Response}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /api/v1/system/attachments [post]
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

// Update godoc
// @Summary      更新附件信息
// @Tags         System / 附件管理
// @Accept       json
// @Produce      json
// @Param        id    path  uint                           true  "附件 ID"
// @Param        body  body  UpdateRequest  true  "附件参数"
// @Success      200  {object}  httpx.Body{data=Response}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /api/v1/system/attachments/{id}/update [post]
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

// UpdateStatus godoc
// @Summary      更新附件状态
// @Tags         System / 附件管理
// @Accept       json
// @Produce      json
// @Param        id    path  uint                                 true  "附件 ID"
// @Param        body  body  UpdateStatusRequest  true  "状态参数"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /api/v1/system/attachments/{id}/status [post]
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
