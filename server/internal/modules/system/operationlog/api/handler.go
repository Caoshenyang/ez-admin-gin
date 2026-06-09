// Package api 提供操作日志模块的 HTTP 请求处理器与路由定义。
package api

import (
	operationlogapp "ez-admin-gin/server/internal/modules/system/operationlog/application"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	httpx "ez-admin-gin/server/internal/pkg/httpx"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	service *operationlogapp.Service
	log     *zap.Logger
}

func NewHandler(service *operationlogapp.Service, log *zap.Logger) *Handler {
	return &Handler{service: service, log: log}
}

// List godoc
// @Summary      查询操作日志列表
// @Tags         System / 操作日志
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "页码"
// @Param        page_size  query     int     false  "每页条数"
// @Param        username   query     string  false  "用户名"
// @Param        method     query     string  false  "请求方法"
// @Param        path       query     string  false  "请求路径"
// @Param        success    query     string  false  "是否成功，支持 true/false/1/0"
// @Success      200  {object}  httpx.Body{data=ListResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/operation-logs [get]
func (h *Handler) List(c *gin.Context) {
	var query ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httpx.Error(c, errorsx.BadRequest("查询参数不正确"), h.log)
		return
	}
	result, err := h.service.List(query)
	if err != nil {
		httpx.WriteError(c, err, "查询操作日志列表失败", h.log)
		return
	}
	httpx.Success(c, result)
}
