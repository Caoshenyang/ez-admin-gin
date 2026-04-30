package config

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
)

var codePattern = regexp.MustCompile(`^[a-z0-9:_-]+$`)

// ListQuery 表示系统配置分页查询参数。
type ListQuery struct {
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
	Keyword   string `form:"keyword"`
	GroupCode string `form:"group_code"`
	Status    int    `form:"status"`
}

// CreateRequest 表示创建配置项的请求体。
type CreateRequest struct {
	GroupCode string                   `json:"group_code"`
	Key       string                   `json:"key"`
	Name      string                   `json:"name"`
	Value     string                   `json:"value"`
	Sort      int                      `json:"sort"`
	Status    model.SystemConfigStatus `json:"status"`
	Remark    string                   `json:"remark"`
}

// UpdateRequest 表示编辑配置项的请求体。
type UpdateRequest struct {
	GroupCode string                   `json:"group_code"`
	Name      string                   `json:"name"`
	Value     string                   `json:"value"`
	Sort      int                      `json:"sort"`
	Status    model.SystemConfigStatus `json:"status"`
	Remark    string                   `json:"remark"`
}

// UpdateStatusRequest 表示单独修改配置状态的请求体。
type UpdateStatusRequest struct {
	Status model.SystemConfigStatus `json:"status"`
}

// Response 表示系统配置对象返回结构。
type Response struct {
	ID        uint                     `json:"id"`
	GroupCode string                   `json:"group_code"`
	Key       string                   `json:"key"`
	Name      string                   `json:"name"`
	Value     string                   `json:"value"`
	Sort      int                      `json:"sort"`
	Status    model.SystemConfigStatus `json:"status"`
	Remark    string                   `json:"remark"`
	CreatedAt time.Time                `json:"created_at"`
	UpdatedAt time.Time                `json:"updated_at"`
}

// ListResponse 表示配置分页结果。
type ListResponse struct {
	Items    []Response `json:"items"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

// ValueResponse 表示按 key 读取配置值的返回结构。
type ValueResponse struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Source string `json:"source"`
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
	groupCode, err := normalizeCode("配置分组", req.GroupCode, 64)
	if err != nil {
		return CreateRequest{}, err
	}
	key, err := normalizeCode("配置键", req.Key, 128)
	if err != nil {
		return CreateRequest{}, err
	}
	name, err := normalizeName(req.Name)
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

	req.GroupCode = groupCode
	req.Key = key
	req.Name = name
	req.Status = status
	req.Remark = remark
	return req, nil
}

// NormalizeUpdateRequest 统一校验并收敛编辑参数。
func NormalizeUpdateRequest(req UpdateRequest) (UpdateRequest, error) {
	groupCode, err := normalizeCode("配置分组", req.GroupCode, 64)
	if err != nil {
		return UpdateRequest{}, err
	}
	name, err := normalizeName(req.Name)
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

	req.GroupCode = groupCode
	req.Name = name
	req.Status = status
	req.Remark = remark
	return req, nil
}

// NormalizeKey 用于按 key 读取配置值时的校验。
func NormalizeKey(key string) (string, error) {
	return normalizeCode("配置键", key, 128)
}

// ValidStatus 判断配置状态是否合法。
func ValidStatus(status model.SystemConfigStatus) bool {
	return status == model.SystemConfigStatusEnabled || status == model.SystemConfigStatusDisabled
}

func normalizeCode(fieldName string, value string, maxLen int) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", apperror.BadRequest(fieldName + "不能为空")
	}
	if len(value) > maxLen {
		return "", apperror.BadRequest(fieldName + "长度不能超过 " + strconv.Itoa(maxLen) + " 个字符")
	}
	if !codePattern.MatchString(value) {
		return "", apperror.BadRequest(fieldName + "只能使用小写字母、数字、冒号、短横线和下划线")
	}

	return value, nil
}

func normalizeName(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", apperror.BadRequest("配置名称不能为空")
	}
	if len(value) > 64 {
		return "", apperror.BadRequest("配置名称不能超过 64 个字符")
	}

	return value, nil
}

func normalizeStatus(status model.SystemConfigStatus, allowDefault bool) (model.SystemConfigStatus, error) {
	if status == 0 && allowDefault {
		status = model.SystemConfigStatusEnabled
	}
	if !ValidStatus(status) {
		return 0, apperror.BadRequest("配置状态不正确")
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

// BuildResponse 把模型对象压成 API 返回结构。
func BuildResponse(item model.SystemConfig) Response {
	return Response{
		ID:        item.ID,
		GroupCode: item.GroupCode,
		Key:       item.ConfigKey,
		Name:      item.Name,
		Value:     item.Value,
		Sort:      item.Sort,
		Status:    item.Status,
		Remark:    item.Remark,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}
