// Package api 提供字典模块的 HTTP 请求处理器与路由定义。
package api

import (
	dictapp "ez-admin-gin/server/internal/modules/system/dict/application"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	httpx "ez-admin-gin/server/internal/pkg/httpx"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	service *dictapp.Service
	log     *zap.Logger
}

func NewHandler(service *dictapp.Service, log *zap.Logger) *Handler {
	return &Handler{service: service, log: log}
}

// ListTypes godoc
// @Summary      查询字典类型列表
// @Tags         System / 字典管理
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "页码"
// @Param        page_size  query     int     false  "每页条数"
// @Param        keyword    query     string  false  "关键词"
// @Param        status     query     int     false  "状态"
// @Success      200  {object}  httpx.Body{data=TypeListResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/dict-types [get]
func (h *Handler) ListTypes(c *gin.Context) {
	var query TypeListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httpx.Error(c, errorsx.BadRequest("查询参数不正确"), h.log)
		return
	}

	result, err := h.service.ListTypes(query)
	if err != nil {
		httpx.WriteError(c, err, "查询字典类型失败", h.log)
		return
	}

	httpx.Success(c, result)
}

// CreateType godoc
// @Summary      创建字典类型
// @Tags         System / 字典管理
// @Accept       json
// @Produce      json
// @Param        body  body  CreateTypeRequest  true  "字典类型参数"
// @Success      200  {object}  httpx.Body{data=TypeResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/dict-types [post]
func (h *Handler) CreateType(c *gin.Context) {
	var req CreateTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.CreateType(req)
	if err != nil {
		httpx.WriteError(c, err, "创建字典类型失败", h.log)
		return
	}

	httpx.Success(c, result)
}

// UpdateType godoc
// @Summary      更新字典类型
// @Tags         System / 字典管理
// @Accept       json
// @Produce      json
// @Param        id    path  uint               true  "字典类型 ID"
// @Param        body  body  UpdateTypeRequest  true  "字典类型参数"
// @Success      200  {object}  httpx.Body{data=TypeResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/dict-types/{id}/update [post]
func (h *Handler) UpdateType(c *gin.Context) {
	typeID, ok := httpx.UintIDParam(c, "id", "字典类型 ID", h.log)
	if !ok {
		return
	}

	var req UpdateTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.UpdateType(typeID, req)
	if err != nil {
		httpx.WriteError(c, err, "更新字典类型失败", h.log)
		return
	}

	httpx.Success(c, result)
}

// UpdateTypeStatus godoc
// @Summary      更新字典类型状态
// @Tags         System / 字典管理
// @Accept       json
// @Produce      json
// @Param        id    path  uint                     true  "字典类型 ID"
// @Param        body  body  UpdateTypeStatusRequest  true  "状态参数"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/dict-types/{id}/status [post]
func (h *Handler) UpdateTypeStatus(c *gin.Context) {
	typeID, ok := httpx.UintIDParam(c, "id", "字典类型 ID", h.log)
	if !ok {
		return
	}

	var req UpdateTypeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	if err := h.service.UpdateTypeStatus(typeID, req.Status); err != nil {
		httpx.WriteError(c, err, "更新字典类型状态失败", h.log)
		return
	}

	httpx.Success(c, gin.H{"id": typeID, "status": req.Status})
}

// DeleteType godoc
// @Summary      删除字典类型
// @Tags         System / 字典管理
// @Accept       json
// @Produce      json
// @Param        id    path  uint  true  "字典类型 ID"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/dict-types/{id}/delete [post]
func (h *Handler) DeleteType(c *gin.Context) {
	typeID, ok := httpx.UintIDParam(c, "id", "字典类型 ID", h.log)
	if !ok {
		return
	}

	if err := h.service.DeleteType(typeID); err != nil {
		httpx.WriteError(c, err, "删除字典类型失败", h.log)
		return
	}

	httpx.Success(c, gin.H{"id": typeID})
}

// ListItems godoc
// @Summary      查询字典项列表
// @Tags         System / 字典管理
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "页码"
// @Param        page_size  query     int     false  "每页条数"
// @Param        type_id    query     uint    true   "字典类型 ID"
// @Param        keyword    query     string  false  "关键词"
// @Param        status     query     int     false  "状态"
// @Success      200  {object}  httpx.Body{data=ItemListResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/dict-items [get]
func (h *Handler) ListItems(c *gin.Context) {
	var query ItemListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httpx.Error(c, errorsx.BadRequest("查询参数不正确"), h.log)
		return
	}

	result, err := h.service.ListItems(query)
	if err != nil {
		httpx.WriteError(c, err, "查询字典项失败", h.log)
		return
	}

	httpx.Success(c, result)
}

// CreateItem godoc
// @Summary      创建字典项
// @Tags         System / 字典管理
// @Accept       json
// @Produce      json
// @Param        body  body  CreateItemRequest  true  "字典项参数"
// @Success      200  {object}  httpx.Body{data=ItemResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/dict-items [post]
func (h *Handler) CreateItem(c *gin.Context) {
	var req CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.CreateItem(req)
	if err != nil {
		httpx.WriteError(c, err, "创建字典项失败", h.log)
		return
	}

	httpx.Success(c, result)
}

// UpdateItem godoc
// @Summary      更新字典项
// @Tags         System / 字典管理
// @Accept       json
// @Produce      json
// @Param        id    path  uint               true  "字典项 ID"
// @Param        body  body  UpdateItemRequest  true  "字典项参数"
// @Success      200  {object}  httpx.Body{data=ItemResponse}
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/dict-items/{id}/update [post]
func (h *Handler) UpdateItem(c *gin.Context) {
	itemID, ok := httpx.UintIDParam(c, "id", "字典项 ID", h.log)
	if !ok {
		return
	}

	var req UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.UpdateItem(itemID, req)
	if err != nil {
		httpx.WriteError(c, err, "更新字典项失败", h.log)
		return
	}

	httpx.Success(c, result)
}

// UpdateItemStatus godoc
// @Summary      更新字典项状态
// @Tags         System / 字典管理
// @Accept       json
// @Produce      json
// @Param        id    path  uint                     true  "字典项 ID"
// @Param        body  body  UpdateItemStatusRequest  true  "状态参数"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/dict-items/{id}/status [post]
func (h *Handler) UpdateItemStatus(c *gin.Context) {
	itemID, ok := httpx.UintIDParam(c, "id", "字典项 ID", h.log)
	if !ok {
		return
	}

	var req UpdateItemStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, errorsx.BadRequest("请求参数不正确"), h.log)
		return
	}

	if err := h.service.UpdateItemStatus(itemID, req.Status); err != nil {
		httpx.WriteError(c, err, "更新字典项状态失败", h.log)
		return
	}

	httpx.Success(c, gin.H{"id": itemID, "status": req.Status})
}

// DeleteItem godoc
// @Summary      删除字典项
// @Tags         System / 字典管理
// @Accept       json
// @Produce      json
// @Param        id    path  uint  true  "字典项 ID"
// @Success      200  {object}  httpx.Body
// @Failure      400  {object}  httpx.Body
// @Failure      401  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /system/dict-items/{id}/delete [post]
func (h *Handler) DeleteItem(c *gin.Context) {
	itemID, ok := httpx.UintIDParam(c, "id", "字典项 ID", h.log)
	if !ok {
		return
	}

	if err := h.service.DeleteItem(itemID); err != nil {
		httpx.WriteError(c, err, "删除字典项失败", h.log)
		return
	}

	httpx.Success(c, gin.H{"id": itemID})
}
