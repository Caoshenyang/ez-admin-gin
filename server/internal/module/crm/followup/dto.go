package followup

import (
	"strconv"
	"strings"
	"time"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
)

// ListQuery 表示客户跟进分页查询参数。
type ListQuery struct {
	Page       int    `form:"page"`
	PageSize   int    `form:"page_size"`
	Keyword    string `form:"keyword"`
	FollowType string `form:"follow_type"`
	CustomerID uint   `form:"customer_id"`
	Status     int    `form:"status"`
}

// CreateRequest 表示创建客户跟进请求体。
type CreateRequest struct {
	CustomerID   uint                         `json:"customer_id"`
	FollowType   string                       `json:"follow_type"`
	Subject      string                       `json:"subject"`
	Content      string                       `json:"content"`
	Result       string                       `json:"result"`
	NextFollowAt *time.Time                   `json:"next_follow_at"`
	Status       model.CustomerFollowUpStatus `json:"status"`
}

// UpdateRequest 表示编辑客户跟进请求体。
type UpdateRequest struct {
	FollowType   string                       `json:"follow_type"`
	Subject      string                       `json:"subject"`
	Content      string                       `json:"content"`
	Result       string                       `json:"result"`
	NextFollowAt *time.Time                   `json:"next_follow_at"`
	Status       model.CustomerFollowUpStatus `json:"status"`
}

// UpdateStatusRequest 表示单独修改客户跟进状态请求体。
type UpdateStatusRequest struct {
	Status model.CustomerFollowUpStatus `json:"status"`
}

// Response 表示客户跟进对象返回结构。
type Response struct {
	ID             uint                         `json:"id"`
	CustomerID     uint                         `json:"customer_id"`
	CustomerName   string                       `json:"customer_name"`
	DepartmentID   uint                         `json:"department_id"`
	DepartmentName string                       `json:"department_name"`
	OwnerUserID    uint                         `json:"owner_user_id"`
	OwnerUsername  string                       `json:"owner_username"`
	OwnerNickname  string                       `json:"owner_nickname"`
	FollowType     string                       `json:"follow_type"`
	Subject        string                       `json:"subject"`
	Content        string                       `json:"content"`
	Result         string                       `json:"result"`
	NextFollowAt   *time.Time                   `json:"next_follow_at"`
	Status         model.CustomerFollowUpStatus `json:"status"`
	CreatedAt      time.Time                    `json:"created_at"`
	UpdatedAt      time.Time                    `json:"updated_at"`
}

// ListResponse 表示客户跟进分页结果。
type ListResponse struct {
	Items    []Response `json:"items"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

// CustomerOption 表示创建跟进时可选的客户。
type CustomerOption struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	DepartmentID   uint   `json:"department_id"`
	DepartmentName string `json:"department_name"`
	OwnerUserID    uint   `json:"owner_user_id"`
	OwnerUsername  string `json:"owner_username"`
	OwnerNickname  string `json:"owner_nickname"`
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
	customerID, err := normalizeCustomerID(req.CustomerID)
	if err != nil {
		return CreateRequest{}, err
	}
	followType, err := normalizeFollowType(req.FollowType)
	if err != nil {
		return CreateRequest{}, err
	}
	subject, err := normalizeSubject(req.Subject)
	if err != nil {
		return CreateRequest{}, err
	}
	content, err := normalizeContent(req.Content)
	if err != nil {
		return CreateRequest{}, err
	}
	result, err := normalizeResult(req.Result)
	if err != nil {
		return CreateRequest{}, err
	}
	status, err := normalizeStatus(req.Status, true)
	if err != nil {
		return CreateRequest{}, err
	}

	req.CustomerID = customerID
	req.FollowType = followType
	req.Subject = subject
	req.Content = content
	req.Result = result
	req.Status = status
	return req, nil
}

// NormalizeUpdateRequest 统一校验并收敛编辑参数。
func NormalizeUpdateRequest(req UpdateRequest) (UpdateRequest, error) {
	followType, err := normalizeFollowType(req.FollowType)
	if err != nil {
		return UpdateRequest{}, err
	}
	subject, err := normalizeSubject(req.Subject)
	if err != nil {
		return UpdateRequest{}, err
	}
	content, err := normalizeContent(req.Content)
	if err != nil {
		return UpdateRequest{}, err
	}
	result, err := normalizeResult(req.Result)
	if err != nil {
		return UpdateRequest{}, err
	}
	status, err := normalizeStatus(req.Status, false)
	if err != nil {
		return UpdateRequest{}, err
	}

	req.FollowType = followType
	req.Subject = subject
	req.Content = content
	req.Result = result
	req.Status = status
	return req, nil
}

// NormalizeStatusFilter 把状态查询参数转换成客户跟进状态。
func NormalizeStatusFilter(value int) (*model.CustomerFollowUpStatus, error) {
	if value == 0 {
		return nil, nil
	}

	status := model.CustomerFollowUpStatus(value)
	if !validStatus(status) {
		return nil, apperror.BadRequest("客户跟进状态不正确")
	}

	return &status, nil
}

// NormalizeCustomerOptionLimit 统一客户选项查询条数。
func NormalizeCustomerOptionLimit(limit int) int {
	if limit < 1 {
		return 50
	}
	if limit > 200 {
		return 200
	}
	return limit
}

// ParseFollowUpID 解析路径参数中的客户跟进 ID。
func ParseFollowUpID(value string) (uint, error) {
	id, err := strconv.ParseUint(strings.TrimSpace(value), 10, 64)
	if err != nil || id == 0 {
		return 0, apperror.BadRequest("客户跟进 ID 不正确")
	}
	return uint(id), nil
}

// BuildResponse 把联合查询视图压成 API 返回结构。
func BuildResponse(item View) Response {
	return Response{
		ID:             item.ID,
		CustomerID:     item.CustomerID,
		CustomerName:   item.CustomerName,
		DepartmentID:   item.DepartmentID,
		DepartmentName: item.DepartmentName,
		OwnerUserID:    item.OwnerUserID,
		OwnerUsername:  item.OwnerUsername,
		OwnerNickname:  item.OwnerNickname,
		FollowType:     item.FollowType,
		Subject:        item.Subject,
		Content:        item.Content,
		Result:         item.Result,
		NextFollowAt:   item.NextFollowAt,
		Status:         item.Status,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
	}
}

func normalizeCustomerID(value uint) (uint, error) {
	if value == 0 {
		return 0, apperror.BadRequest("客户 ID 不正确")
	}
	return value, nil
}

func normalizeFollowType(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", apperror.BadRequest("跟进方式不能为空")
	}
	if len(value) > 32 {
		return "", apperror.BadRequest("跟进方式不能超过 32 个字符")
	}
	return value, nil
}

func normalizeSubject(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", apperror.BadRequest("跟进主题不能为空")
	}
	if len(value) > 128 {
		return "", apperror.BadRequest("跟进主题不能超过 128 个字符")
	}
	return value, nil
}

func normalizeContent(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", apperror.BadRequest("跟进内容不能为空")
	}
	if len(value) > 1000 {
		return "", apperror.BadRequest("跟进内容不能超过 1000 个字符")
	}
	return value, nil
}

func normalizeResult(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) > 255 {
		return "", apperror.BadRequest("跟进结果不能超过 255 个字符")
	}
	return value, nil
}

func normalizeStatus(status model.CustomerFollowUpStatus, allowDefault bool) (model.CustomerFollowUpStatus, error) {
	if status == 0 && allowDefault {
		status = model.CustomerFollowUpStatusPending
	}
	if !validStatus(status) {
		return 0, apperror.BadRequest("客户跟进状态不正确")
	}
	return status, nil
}

func validStatus(status model.CustomerFollowUpStatus) bool {
	return status == model.CustomerFollowUpStatusPending ||
		status == model.CustomerFollowUpStatusCompleted ||
		status == model.CustomerFollowUpStatusClosed
}
