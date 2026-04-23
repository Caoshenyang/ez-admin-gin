package system

import (
	"strings"
	"time"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// OperationLogHandler 负责操作日志查询接口。
type OperationLogHandler struct {
	db  *gorm.DB
	log *zap.Logger
}

// NewOperationLogHandler 创建操作日志 Handler。
func NewOperationLogHandler(db *gorm.DB, log *zap.Logger) *OperationLogHandler {
	return &OperationLogHandler{
		db:  db,
		log: log,
	}
}

type operationLogListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Username string `form:"username"`
	Method   string `form:"method"`
	Path     string `form:"path"`
	Success  string `form:"success"`
}

type operationLogResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	Username     string    `json:"username"`
	Method       string    `json:"method"`
	Path         string    `json:"path"`
	RoutePath    string    `json:"route_path"`
	Query        string    `json:"query"`
	IP           string    `json:"ip"`
	UserAgent    string    `json:"user_agent"`
	StatusCode   int       `json:"status_code"`
	LatencyMs    int64     `json:"latency_ms"`
	Success      bool      `json:"success"`
	ErrorMessage string    `json:"error_message"`
	CreatedAt    time.Time `json:"created_at"`
}

type operationLogListResponse struct {
	Items    []operationLogResponse `json:"items"`
	Total    int64                  `json:"total"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
}

// List 返回操作日志分页列表。
func (h *OperationLogHandler) List(c *gin.Context) {
	var query operationLogListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, apperror.BadRequest("查询参数不正确"), h.log)
		return
	}

	page, pageSize := normalizeOperationLogPage(query.Page, query.PageSize)
	queryDB := h.db.Model(&model.OperationLog{})

	username := strings.TrimSpace(query.Username)
	if username != "" {
		queryDB = queryDB.Where("username = ?", username)
	}

	method := strings.ToUpper(strings.TrimSpace(query.Method))
	if method != "" {
		queryDB = queryDB.Where("method = ?", method)
	}

	path := strings.TrimSpace(query.Path)
	if path != "" {
		queryDB = queryDB.Where("path LIKE ?", "%"+path+"%")
	}

	if query.Success != "" {
		success, ok := parseOperationLogSuccess(query.Success)
		if !ok {
			response.Error(c, apperror.BadRequest("成功状态不正确"), h.log)
			return
		}
		queryDB = queryDB.Where("success = ?", success)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		response.Error(c, apperror.Internal("查询操作日志总数失败", err), h.log)
		return
	}

	var logs []model.OperationLog
	if err := queryDB.
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&logs).Error; err != nil {
		response.Error(c, apperror.Internal("查询操作日志列表失败", err), h.log)
		return
	}

	items := make([]operationLogResponse, 0, len(logs))
	for _, item := range logs {
		items = append(items, buildOperationLogResponse(item))
	}

	response.Success(c, operationLogListResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func normalizeOperationLogPage(page int, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return page, pageSize
}

func parseOperationLogSuccess(value string) (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true", "1":
		return true, true
	case "false", "0":
		return false, true
	default:
		return false, false
	}
}

func buildOperationLogResponse(item model.OperationLog) operationLogResponse {
	return operationLogResponse{
		ID:           item.ID,
		UserID:       item.UserID,
		Username:     item.Username,
		Method:       item.Method,
		Path:         item.Path,
		RoutePath:    item.RoutePath,
		Query:        item.Query,
		IP:           item.IP,
		UserAgent:    item.UserAgent,
		StatusCode:   item.StatusCode,
		LatencyMs:    item.LatencyMs,
		Success:      item.Success,
		ErrorMessage: item.ErrorMessage,
		CreatedAt:    item.CreatedAt,
	}
}
