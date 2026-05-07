package api

import (
	loginlogapp "ez-admin-gin/server/internal/modules/system/loginlog/application"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	httpx "ez-admin-gin/server/internal/pkg/httpx"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	service *loginlogapp.Service
	log     *zap.Logger
}

func NewHandler(service *loginlogapp.Service, log *zap.Logger) *Handler {
	return &Handler{service: service, log: log}
}

// List godoc
// @Summary      查询登录日志列表
// @Tags         System / 登录日志
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "页码"
// @Param        page_size  query     int     false  "每页条数"
// @Param        username   query     string  false  "用户名"
// @Param        ip         query     string  false  "IP 地址"
// @Param        status     query     int     false  "登录状态"
// @Success      200  {object}  httpx.Body{data=ListResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/login-logs [get]
func (h *Handler) List(c *gin.Context) {
	var query ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httpx.Error(c, errorsx.BadRequest("查询参数不正确"), h.log)
		return
	}
	result, err := h.service.List(query)
	if err != nil {
		httpx.WriteError(c, err, "查询登录日志列表失败", h.log)
		return
	}
	httpx.Success(c, result)
}
