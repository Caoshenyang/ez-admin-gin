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

// LoginLogHandler 负责登录日志查询接口。
type LoginLogHandler struct {
	db  *gorm.DB
	log *zap.Logger
}

// NewLoginLogHandler 创建登录日志 Handler。
func NewLoginLogHandler(db *gorm.DB, log *zap.Logger) *LoginLogHandler {
	return &LoginLogHandler{
		db:  db,
		log: log,
	}
}

type loginLogListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Username string `form:"username"`
	IP       string `form:"ip"`
	Status   int    `form:"status"`
}

type loginLogResponse struct {
	ID        uint                 `json:"id"`
	UserID    uint                 `json:"user_id"`
	Username  string               `json:"username"`
	Status    model.LoginLogStatus `json:"status"`
	Message   string               `json:"message"`
	IP        string               `json:"ip"`
	UserAgent string               `json:"user_agent"`
	CreatedAt time.Time            `json:"created_at"`
}

type loginLogListResponse struct {
	Items    []loginLogResponse `json:"items"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}

// List 返回登录日志分页列表。
func (h *LoginLogHandler) List(c *gin.Context) {
	var query loginLogListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, apperror.BadRequest("查询参数不正确"), h.log)
		return
	}

	page, pageSize := normalizeLoginLogPage(query.Page, query.PageSize)
	queryDB := h.db.Model(&model.LoginLog{})

	username := strings.TrimSpace(query.Username)
	if username != "" {
		queryDB = queryDB.Where("username = ?", username)
	}

	ip := strings.TrimSpace(query.IP)
	if ip != "" {
		queryDB = queryDB.Where("ip = ?", ip)
	}

	if query.Status != 0 {
		status := model.LoginLogStatus(query.Status)
		if !validLoginLogStatus(status) {
			response.Error(c, apperror.BadRequest("登录状态不正确"), h.log)
			return
		}
		queryDB = queryDB.Where("status = ?", status)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		response.Error(c, apperror.Internal("查询登录日志总数失败", err), h.log)
		return
	}

	var logs []model.LoginLog
	if err := queryDB.
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&logs).Error; err != nil {
		response.Error(c, apperror.Internal("查询登录日志列表失败", err), h.log)
		return
	}

	items := make([]loginLogResponse, 0, len(logs))
	for _, item := range logs {
		items = append(items, buildLoginLogResponse(item))
	}

	response.Success(c, loginLogListResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func normalizeLoginLogPage(page int, pageSize int) (int, int) {
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

func validLoginLogStatus(status model.LoginLogStatus) bool {
	return status == model.LoginLogStatusSuccess || status == model.LoginLogStatusFailed
}

func buildLoginLogResponse(item model.LoginLog) loginLogResponse {
	return loginLogResponse{
		ID:        item.ID,
		UserID:    item.UserID,
		Username:  item.Username,
		Status:    item.Status,
		Message:   item.Message,
		IP:        item.IP,
		UserAgent: item.UserAgent,
		CreatedAt: item.CreatedAt,
	}
}
