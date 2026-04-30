package notice

import (
	"strconv"
	"strings"
	"time"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
)

// ListQuery 表示公告分页查询参数。
type ListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
	Status   int    `form:"status"`
}

// CreateRequest 表示创建公告请求体。
type CreateRequest struct {
	Title   string             `json:"title"`
	Content string             `json:"content"`
	Sort    int                `json:"sort"`
	Status  model.NoticeStatus `json:"status"`
	Remark  string             `json:"remark"`
}

// UpdateRequest 表示编辑公告请求体。
type UpdateRequest struct {
	Title   string             `json:"title"`
	Content string             `json:"content"`
	Sort    int                `json:"sort"`
	Status  model.NoticeStatus `json:"status"`
	Remark  string             `json:"remark"`
}

// UpdateStatusRequest 表示单独修改公告状态的请求体。
type UpdateStatusRequest struct {
	Status model.NoticeStatus `json:"status"`
}

// Response 表示公告对象返回结构。
type Response struct {
	ID        uint               `json:"id"`
	Title     string             `json:"title"`
	Content   string             `json:"content"`
	Sort      int                `json:"sort"`
	Status    model.NoticeStatus `json:"status"`
	Remark    string             `json:"remark"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// ListResponse 表示公告分页结果。
type ListResponse struct {
	Items    []Response `json:"items"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

// NormalizePage 统一分页边界。
func NormalizePage(page int, pageSize int) (int, int) {
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

// NormalizeCreateRequest 统一校验并收敛创建参数。
func NormalizeCreateRequest(req CreateRequest) (CreateRequest, error) {
	title, err := normalizeTitle(req.Title)
	if err != nil {
		return CreateRequest{}, err
	}
	status, err := normalizeStatus(req.Status, true)
	if err != nil {
		return CreateRequest{}, err
	}
	remark, err := normalizeRemark(req.Remark)
	if err != nil {
		return CreateRequest{}, err
	}

	req.Title = title
	req.Content = strings.TrimSpace(req.Content)
	req.Status = status
	req.Remark = remark
	return req, nil
}

// NormalizeUpdateRequest 统一校验并收敛编辑参数。
func NormalizeUpdateRequest(req UpdateRequest) (UpdateRequest, error) {
	title, err := normalizeTitle(req.Title)
	if err != nil {
		return UpdateRequest{}, err
	}
	status, err := normalizeStatus(req.Status, false)
	if err != nil {
		return UpdateRequest{}, err
	}
	remark, err := normalizeRemark(req.Remark)
	if err != nil {
		return UpdateRequest{}, err
	}

	req.Title = title
	req.Content = strings.TrimSpace(req.Content)
	req.Status = status
	req.Remark = remark
	return req, nil
}

func normalizeTitle(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", apperror.BadRequest("公告标题不能为空")
	}
	if len(value) > 128 {
		return "", apperror.BadRequest("公告标题不能超过 128 个字符")
	}
	return value, nil
}

func normalizeStatus(status model.NoticeStatus, allowDefault bool) (model.NoticeStatus, error) {
	if status == 0 && allowDefault {
		status = model.NoticeStatusEnabled
	}
	if status != model.NoticeStatusEnabled && status != model.NoticeStatusDisabled {
		return 0, apperror.BadRequest("公告状态不正确")
	}
	return status, nil
}

func normalizeRemark(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) > 255 {
		return "", apperror.BadRequest("备注不能超过 255 个字符")
	}
	return value, nil
}

// ParseNoticeID 解析路径参数中的公告 ID。
func ParseNoticeID(value string) (uint, error) {
	id, err := strconv.ParseUint(strings.TrimSpace(value), 10, 64)
	if err != nil || id == 0 {
		return 0, apperror.BadRequest("公告 ID 不正确")
	}
	return uint(id), nil
}

// NormalizeStatusFilter 把状态查询参数转换成公告状态。
func NormalizeStatusFilter(value int) (*model.NoticeStatus, error) {
	if value == 0 {
		return nil, nil
	}

	status := model.NoticeStatus(value)
	if status != model.NoticeStatusEnabled && status != model.NoticeStatusDisabled {
		return nil, apperror.BadRequest("公告状态不正确")
	}

	return &status, nil
}

// BuildResponse 把模型对象压成 API 返回结构。
func BuildResponse(item model.Notice) Response {
	return Response{
		ID:        item.ID,
		Title:     item.Title,
		Content:   item.Content,
		Sort:      item.Sort,
		Status:    item.Status,
		Remark:    item.Remark,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}
