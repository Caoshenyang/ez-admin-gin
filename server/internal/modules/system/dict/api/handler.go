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
