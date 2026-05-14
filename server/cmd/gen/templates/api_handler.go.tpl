package api

import (
	{{ .Module }}app "ez-admin-gin/server/internal/modules/{{ .Group }}/{{ .Module }}/application"
	{{ .Module }}domain "ez-admin-gin/server/internal/modules/{{ .Group }}/{{ .Module }}/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	httpx "ez-admin-gin/server/internal/pkg/httpx"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	service *{{ .Module }}app.Service
	log     *zap.Logger
}

func NewHandler(service *{{ .Module }}app.Service, log *zap.Logger) *Handler {
	return &Handler{service: service, log: log}
}

// List godoc
// @Summary      查询{{ .Label }}列表
// @Tags         {{ title .Group }} / {{ .Label }}管理
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "页码"
// @Param        page_size  query     int     false  "每页条数"
// @Success      200  {object}  httpx.Body{data={{ .Module }}domain.ListResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /{{ .Group }}/{{ .Module }}s [get]
func (h *Handler) List(c *gin.Context) {
	var query ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httpx.Error(c, errorsx.BadRequest("查询参数不正确"), h.log)
		return
	}
	result, err := h.service.List(query)
	if err != nil {
		httpx.WriteError(c, err, "查询{{ .Label }}列表失败", h.log)
		return
	}
	httpx.Success(c, result)
}

// Create godoc
// @Summary      创建{{ .Label }}
// @Tags         {{ title .Group }} / {{ .Label }}管理
// @Accept       json
// @Produce      json
// @Param        body  body  {{ .Module }}domain.CreateRequest  true  "{{ .Label }}参数"
// @Success      200  {object}  httpx.Body{data={{ .Module }}domain.Response}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /{{ .Group }}/{{ .Module }}s [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}
	result, err := h.service.Create(req)
	if err != nil {
		httpx.WriteError(c, err, "创建{{ .Label }}失败", h.log)
		return
	}
	httpx.Success(c, result)
}

// Update godoc
// @Summary      更新{{ .Label }}
// @Tags         {{ title .Group }} / {{ .Label }}管理
// @Accept       json
// @Produce      json
// @Param        id    path  uint                          true  "{{ .Label }} ID"
// @Param        body  body  {{ .Module }}domain.UpdateRequest  true  "{{ .Label }}参数"
// @Success      200  {object}  httpx.Body{data={{ .Module }}domain.Response}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /{{ .Group }}/{{ .Module }}s/{id}/update [post]
func (h *Handler) Update(c *gin.Context) {
	id, ok := httpx.UintIDParam(c, "id", "{{ .Label }} ID", h.log)
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
		httpx.WriteError(c, err, "更新{{ .Label }}失败", h.log)
		return
	}
	httpx.Success(c, result)
}

{{- if .HasStatus }}

// UpdateStatus godoc
// @Summary      更新{{ .Label }}状态
// @Tags         {{ title .Group }} / {{ .Label }}管理
// @Accept       json
// @Produce      json
// @Param        id    path  uint                              true  "{{ .Label }} ID"
// @Param        body  body  {{ .Module }}domain.UpdateStatusRequest  true  "状态参数"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /{{ .Group }}/{{ .Module }}s/{id}/status [post]
func (h *Handler) UpdateStatus(c *gin.Context) {
	id, ok := httpx.UintIDParam(c, "id", "{{ .Label }} ID", h.log)
	if !ok {
		return
	}
	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}
{{- $status := .StatusField }}
{{- if $status }}
	if err := h.service.UpdateStatus(id, req.{{ title $status.Name }}); err != nil {
		httpx.WriteError(c, err, "更新{{ $.Label }}状态失败", h.log)
		return
	}
	httpx.Success(c, gin.H{"id": id, "{{ $status.Name }}": req.{{ title $status.Name }}})
{{- end }}
}
{{- end }}

var _ = {{ .Module }}domain.PermissionList
