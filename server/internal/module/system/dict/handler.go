package dict

import (
	"errors"
	"strconv"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 负责字典模块的 HTTP 协议层绑定与输出。
type Handler struct {
	service *Service
	log     *zap.Logger
}

// NewHandler 创建字典 Handler。
func NewHandler(service *Service, log *zap.Logger) *Handler {
	return &Handler{service: service, log: log}
}

// ListTypes 返回字典类型分页列表。
func (h *Handler) ListTypes(c *gin.Context) {
	var query TypeListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, apperror.BadRequest("查询参数不正确"), h.log)
		return
	}

	result, err := h.service.ListTypes(query)
	if err != nil {
		writeError(c, err, "查询字典类型失败", h.log)
		return
	}

	response.Success(c, result)
}

// CreateType 创建字典类型。
func (h *Handler) CreateType(c *gin.Context) {
	var req CreateTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.CreateType(req)
	if err != nil {
		writeError(c, err, "创建字典类型失败", h.log)
		return
	}

	response.Success(c, result)
}

// UpdateType 编辑字典类型。
func (h *Handler) UpdateType(c *gin.Context) {
	typeID, ok := uintIDParam(c, "字典类型 ID 不正确", h.log)
	if !ok {
		return
	}

	var req UpdateTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.UpdateType(typeID, req)
	if err != nil {
		writeError(c, err, "更新字典类型失败", h.log)
		return
	}

	response.Success(c, result)
}

// UpdateTypeStatus 修改字典类型状态。
func (h *Handler) UpdateTypeStatus(c *gin.Context) {
	typeID, ok := uintIDParam(c, "字典类型 ID 不正确", h.log)
	if !ok {
		return
	}

	var req UpdateTypeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	if err := h.service.UpdateTypeStatus(typeID, req.Status); err != nil {
		writeError(c, err, "更新字典类型状态失败", h.log)
		return
	}

	response.Success(c, gin.H{"id": typeID, "status": req.Status})
}

// ListItems 返回字典项分页列表。
func (h *Handler) ListItems(c *gin.Context) {
	var query ItemListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, apperror.BadRequest("查询参数不正确"), h.log)
		return
	}

	result, err := h.service.ListItems(query)
	if err != nil {
		writeError(c, err, "查询字典项失败", h.log)
		return
	}

	response.Success(c, result)
}

// CreateItem 创建字典项。
func (h *Handler) CreateItem(c *gin.Context) {
	var req CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.CreateItem(req)
	if err != nil {
		writeError(c, err, "创建字典项失败", h.log)
		return
	}

	response.Success(c, result)
}

// UpdateItem 编辑字典项。
func (h *Handler) UpdateItem(c *gin.Context) {
	itemID, ok := uintIDParam(c, "字典项 ID 不正确", h.log)
	if !ok {
		return
	}

	var req UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	result, err := h.service.UpdateItem(itemID, req)
	if err != nil {
		writeError(c, err, "更新字典项失败", h.log)
		return
	}

	response.Success(c, result)
}

// UpdateItemStatus 修改字典项状态。
func (h *Handler) UpdateItemStatus(c *gin.Context) {
	itemID, ok := uintIDParam(c, "字典项 ID 不正确", h.log)
	if !ok {
		return
	}

	var req UpdateItemStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	if err := h.service.UpdateItemStatus(itemID, req.Status); err != nil {
		writeError(c, err, "更新字典项状态失败", h.log)
		return
	}

	response.Success(c, gin.H{"id": itemID, "status": req.Status})
}

func uintIDParam(c *gin.Context, message string, log *zap.Logger) (uint, bool) {
	rawID := c.Param("id")
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 {
		response.Error(c, apperror.BadRequest(message), log)
		return 0, false
	}
	return uint(id), true
}

func writeError(c *gin.Context, err error, fallbackMessage string, log *zap.Logger) {
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		response.Error(c, appErr, log)
		return
	}

	response.Error(c, apperror.Internal(fallbackMessage, err), log)
}
