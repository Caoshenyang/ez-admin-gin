package domain

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/model"
)

var codePattern = regexp.MustCompile(`^[a-z0-9:_-]+$`)

type ListQuery struct {
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
	Keyword   string `form:"keyword"`
	GroupCode string `form:"group_code"`
	Status    int    `form:"status"`
}

type CreateRequest struct {
	GroupCode string                   `json:"group_code"`
	Key       string                   `json:"key"`
	Name      string                   `json:"name"`
	Value     string                   `json:"value"`
	Sort      int                      `json:"sort"`
	Status    model.SystemConfigStatus `json:"status"`
	Remark    string                   `json:"remark"`
}

type UpdateRequest struct {
	GroupCode string                   `json:"group_code"`
	Name      string                   `json:"name"`
	Value     string                   `json:"value"`
	Sort      int                      `json:"sort"`
	Status    model.SystemConfigStatus `json:"status"`
	Remark    string                   `json:"remark"`
}

type UpdateStatusRequest struct {
	Status model.SystemConfigStatus `json:"status"`
}

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

type ListResponse struct {
	Items    []Response `json:"items"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

type ValueResponse struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Source string `json:"source"`
}

type Entity = model.SystemConfig

const (
	PermissionList         = "system:config:list"
	PermissionCreate       = "system:config:create"
	PermissionUpdate       = "system:config:update"
	PermissionUpdateStatus = "system:config:update_status"
)

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

func NormalizeKey(key string) (string, error) {
	return normalizeCode("配置键", key, 128)
}

func ValidStatus(status model.SystemConfigStatus) bool {
	return status == model.SystemConfigStatusEnabled || status == model.SystemConfigStatusDisabled
}

func normalizeCode(fieldName string, value string, maxLen int) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errorsx.BadRequest(fieldName + "不能为空")
	}
	if len(value) > maxLen {
		return "", errorsx.BadRequest(fieldName + "长度不能超过 " + strconv.Itoa(maxLen) + " 个字符")
	}
	if !codePattern.MatchString(value) {
		return "", errorsx.BadRequest(fieldName + "只能使用小写字母、数字、冒号、短横线和下划线")
	}
	return value, nil
}

func normalizeName(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errorsx.BadRequest("配置名称不能为空")
	}
	if len(value) > 64 {
		return "", errorsx.BadRequest("配置名称不能超过 64 个字符")
	}
	return value, nil
}

func normalizeStatus(status model.SystemConfigStatus, allowDefault bool) (model.SystemConfigStatus, error) {
	if status == 0 && allowDefault {
		status = model.SystemConfigStatusEnabled
	}
	if !ValidStatus(status) {
		return 0, errorsx.BadRequest("配置状态不正确")
	}
	return status, nil
}

func normalizeRemark(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) > 255 {
		return "", errorsx.BadRequest("备注不能超过 255 个字符")
	}
	return value, nil
}

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
