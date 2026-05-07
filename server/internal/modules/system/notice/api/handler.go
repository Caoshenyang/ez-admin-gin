package api

import (
	noticeapp "ez-admin-gin/server/internal/modules/system/notice/application"
	noticedomain "ez-admin-gin/server/internal/modules/system/notice/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	httpx "ez-admin-gin/server/internal/pkg/httpx"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	service *noticeapp.Service
	log     *zap.Logger
}

func NewHandler(service *noticeapp.Service, log *zap.Logger) *Handler {
	return &Handler{service: service, log: log}
}

// List godoc
// @Summary      查询公告列表
// @Tags         System / 公告管理
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "页码"
// @Param        page_size  query     int     false  "每页条数"
// @Param        keyword    query     string  false  "关键词"
// @Param        status     query     int     false  "状态"
// @Success      200  {object}  httpx.Body{data=noticedomain.ListResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /api/v1/system/notices [get]
func (h *Handler) List(c *gin.Context) {
	var query ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httpx.Error(c, errorsx.BadRequest("查询参数不正确"), h.log)
		return
	}
	result, err := h.service.List(query)
	if err != nil {
		httpx.WriteError(c, err, "查询公告列表失败", h.log)
		return
	}
	httpx.Success(c, result)
}

// Create godoc
// @Summary      创建公告
// @Tags         System / 公告管理
// @Accept       json
// @Produce      json
// @Param        body  body  noticedomain.CreateRequest  true  "公告参数"
// @Success      200  {object}  httpx.Body{data=noticedomain.Response}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /api/v1/system/notices [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}
	result, err := h.service.Create(req)
	if err != nil {
		httpx.WriteError(c, err, "创建公告失败", h.log)
		return
	}
	httpx.Success(c, result)
}

// Update godoc
// @Summary      更新公告
// @Tags         System / 公告管理
// @Accept       json
// @Produce      json
// @Param        id    path  uint                        true  "公告 ID"
// @Param        body  body  noticedomain.UpdateRequest  true  "公告参数"
// @Success      200  {object}  httpx.Body{data=noticedomain.Response}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /api/v1/system/notices/{id}/update [post]
func (h *Handler) Update(c *gin.Context) {
	noticeID, ok := httpx.UintIDParam(c, "id", "公告 ID", h.log)
	if !ok {
		return
	}
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}
	result, err := h.service.Update(noticeID, req)
	if err != nil {
		httpx.WriteError(c, err, "更新公告失败", h.log)
		return
	}
	httpx.Success(c, result)
}

// UpdateStatus godoc
// @Summary      更新公告状态
// @Tags         System / 公告管理
// @Accept       json
// @Produce      json
// @Param        id    path  uint                              true  "公告 ID"
// @Param        body  body  noticedomain.UpdateStatusRequest  true  "状态参数"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /api/v1/system/notices/{id}/status [post]
func (h *Handler) UpdateStatus(c *gin.Context) {
	noticeID, ok := httpx.UintIDParam(c, "id", "公告 ID", h.log)
	if !ok {
		return
	}
	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}
	if err := h.service.UpdateStatus(noticeID, req.Status); err != nil {
		httpx.WriteError(c, err, "更新公告状态失败", h.log)
		return
	}
	httpx.Success(c, gin.H{"id": noticeID, "status": req.Status})
}

var _ = noticedomain.PermissionList
