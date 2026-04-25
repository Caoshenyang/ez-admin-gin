package system

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// NoticeHandler 负责公告管理接口。
type NoticeHandler struct {
	db  *gorm.DB
	log *zap.Logger
}

// NewNoticeHandler 创建公告 Handler。
func NewNoticeHandler(db *gorm.DB, log *zap.Logger) *NoticeHandler {
	return &NoticeHandler{db: db, log: log}
}

type noticeListQuery struct {
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
	Keyword   string `form:"keyword"`
	Status    int    `form:"status"`
}

type createNoticeRequest struct {
	Title   string          `json:"title"`
	Content string          `json:"content"`
	Sort    int             `json:"sort"`
	Status  model.NoticeStatus `json:"status"`
	Remark  string          `json:"remark"`
}

type updateNoticeRequest struct {
	Title   string          `json:"title"`
	Content string          `json:"content"`
	Sort    int             `json:"sort"`
	Status  model.NoticeStatus `json:"status"`
	Remark  string          `json:"remark"`
}

type updateNoticeStatusRequest struct {
	Status model.NoticeStatus `json:"status"`
}

type noticeResponse struct {
	ID        uint             `json:"id"`
	Title     string           `json:"title"`
	Content   string           `json:"content"`
	Sort      int              `json:"sort"`
	Status    model.NoticeStatus `json:"status"`
	Remark    string           `json:"remark"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type noticeListResponse struct {
	Items    []noticeResponse `json:"items"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
}

// List 返回公告分页列表。
func (h *NoticeHandler) List(c *gin.Context) {
	var query noticeListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, apperror.BadRequest("查询参数不正确"), h.log)
		return
	}

	page, pageSize := normalizeNoticePage(query.Page, query.PageSize)
	queryDB := h.db.Model(&model.Notice{})

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where("title LIKE ?", like)
	}

	if query.Status != 0 {
		status := model.NoticeStatus(query.Status)
		if !validNoticeStatus(status) {
			response.Error(c, apperror.BadRequest("公告状态不正确"), h.log)
			return
		}
		queryDB = queryDB.Where("status = ?", status)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		response.Error(c, apperror.Internal("查询公告总数失败", err), h.log)
		return
	}

	var notices []model.Notice
	if err := queryDB.
		Order("sort ASC, id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&notices).Error; err != nil {
		response.Error(c, apperror.Internal("查询公告列表失败", err), h.log)
		return
	}

	items := make([]noticeResponse, 0, len(notices))
	for _, n := range notices {
		items = append(items, buildNoticeResponse(n))
	}

	response.Success(c, noticeListResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// Create 创建公告。
func (h *NoticeHandler) Create(c *gin.Context) {
	var req createNoticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		response.Error(c, apperror.BadRequest("公告标题不能为空"), h.log)
		return
	}

	if len(title) > 128 {
		response.Error(c, apperror.BadRequest("公告标题不能超过 128 个字符"), h.log)
		return
	}

	status := req.Status
	if status == 0 {
		status = model.NoticeStatusEnabled
	}
	if !validNoticeStatus(status) {
		response.Error(c, apperror.BadRequest("公告状态不正确"), h.log)
		return
	}

	notice := model.Notice{
		Title:   title,
		Content: req.Content,
		Sort:    req.Sort,
		Status:  status,
		Remark:  strings.TrimSpace(req.Remark),
	}

	if err := h.db.Create(&notice).Error; err != nil {
		response.Error(c, apperror.Internal("创建公告失败", err), h.log)
		return
	}

	response.Success(c, buildNoticeResponse(notice))
}

// Update 编辑公告。
func (h *NoticeHandler) Update(c *gin.Context) {
	noticeID, ok := noticeIDParam(c, h.log)
	if !ok {
		return
	}

	var req updateNoticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		response.Error(c, apperror.BadRequest("公告标题不能为空"), h.log)
		return
	}

	if !validNoticeStatus(req.Status) {
		response.Error(c, apperror.BadRequest("公告状态不正确"), h.log)
		return
	}

	var notice model.Notice
	if err := h.db.First(&notice, noticeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, apperror.NotFound("公告不存在"), h.log)
			return
		}
		response.Error(c, apperror.Internal("查询公告失败", err), h.log)
		return
	}

	if err := h.db.Model(&notice).Updates(map[string]any{
		"title":   title,
		"content": req.Content,
		"sort":    req.Sort,
		"status":  req.Status,
		"remark":  strings.TrimSpace(req.Remark),
	}).Error; err != nil {
		response.Error(c, apperror.Internal("更新公告失败", err), h.log)
		return
	}

	notice.Title = title
	notice.Content = req.Content
	notice.Sort = req.Sort
	notice.Status = req.Status
	notice.Remark = strings.TrimSpace(req.Remark)

	response.Success(c, buildNoticeResponse(notice))
}

// UpdateStatus 修改公告状态。
func (h *NoticeHandler) UpdateStatus(c *gin.Context) {
	noticeID, ok := noticeIDParam(c, h.log)
	if !ok {
		return
	}

	var req updateNoticeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	if !validNoticeStatus(req.Status) {
		response.Error(c, apperror.BadRequest("公告状态不正确"), h.log)
		return
	}

	var notice model.Notice
	if err := h.db.First(&notice, noticeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, apperror.NotFound("公告不存在"), h.log)
			return
		}
		response.Error(c, apperror.Internal("查询公告失败", err), h.log)
		return
	}

	if err := h.db.Model(&notice).Update("status", req.Status).Error; err != nil {
		response.Error(c, apperror.Internal("更新公告状态失败", err), h.log)
		return
	}

	response.Success(c, gin.H{
		"id":     noticeID,
		"status": req.Status,
	})
}

func normalizeNoticePage(page int, pageSize int) (int, int) {
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

func validNoticeStatus(status model.NoticeStatus) bool {
	return status == model.NoticeStatusEnabled || status == model.NoticeStatusDisabled
}

func noticeIDParam(c *gin.Context, log *zap.Logger) (uint, bool) {
	rawID := c.Param("id")
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 {
		response.Error(c, apperror.BadRequest("公告 ID 不正确"), log)
		return 0, false
	}
	return uint(id), true
}

func buildNoticeResponse(n model.Notice) noticeResponse {
	return noticeResponse{
		ID:        n.ID,
		Title:     n.Title,
		Content:   n.Content,
		Sort:      n.Sort,
		Status:    n.Status,
		Remark:    n.Remark,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}
