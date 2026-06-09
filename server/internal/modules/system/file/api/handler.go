// Package api 提供文件上传模块的 HTTP 请求处理器与路由定义。
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

// List godoc
// @Summary      查询文件列表
// @Tags         System / 文件管理
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "页码"
// @Param        page_size  query     int     false  "每页条数"
// @Param        keyword    query     string  false  "关键词"
// @Param        ext        query     string  false  "文件后缀"
// @Param        status     query     int     false  "状态"
// @Success      200  {object}  httpx.Body{data=ListResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/files [get]
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

// Upload godoc
// @Summary      上传文件
// @Tags         System / 文件管理
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "文件"
// @Success      200  {object}  httpx.Body{data=Response}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/files [post]
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
