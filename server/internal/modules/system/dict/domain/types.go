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

type TypeListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
	Status   int    `form:"status"`
}

type CreateTypeRequest struct {
	Code   string                 `json:"code"`
	Name   string                 `json:"name"`
	Sort   int                    `json:"sort"`
	Status model.SystemDictStatus `json:"status"`
	Remark string                 `json:"remark"`
}

type UpdateTypeRequest struct {
	Name   string                 `json:"name"`
	Sort   int                    `json:"sort"`
	Status model.SystemDictStatus `json:"status"`
	Remark string                 `json:"remark"`
}

type UpdateTypeStatusRequest struct {
	Status model.SystemDictStatus `json:"status"`
}

type TypeResponse struct {
	ID        uint                   `json:"id"`
	Code      string                 `json:"code"`
	Name      string                 `json:"name"`
	Sort      int                    `json:"sort"`
	Status    model.SystemDictStatus `json:"status"`
	Remark    string                 `json:"remark"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type TypeListResponse struct {
	Items    []TypeResponse `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

type ItemListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	TypeID   uint   `form:"type_id"`
	Keyword  string `form:"keyword"`
	Status   int    `form:"status"`
}

type CreateItemRequest struct {
	TypeID  uint                   `json:"type_id"`
	ItemKey string                 `json:"item_key"`
	Label   string                 `json:"label"`
	Value   string                 `json:"value"`
	TagType string                 `json:"tag_type"`
	Sort    int                    `json:"sort"`
	Status  model.SystemDictStatus `json:"status"`
	Remark  string                 `json:"remark"`
}

type UpdateItemRequest struct {
	Label   string                 `json:"label"`
	Value   string                 `json:"value"`
	TagType string                 `json:"tag_type"`
	Sort    int                    `json:"sort"`
	Status  model.SystemDictStatus `json:"status"`
	Remark  string                 `json:"remark"`
}

type UpdateItemStatusRequest struct {
	Status model.SystemDictStatus `json:"status"`
}

type ItemResponse struct {
	ID        uint                   `json:"id"`
	TypeID    uint                   `json:"type_id"`
	ItemKey   string                 `json:"item_key"`
	Label     string                 `json:"label"`
	Value     string                 `json:"value"`
	TagType   string                 `json:"tag_type"`
	Sort      int                    `json:"sort"`
	Status    model.SystemDictStatus `json:"status"`
	Remark    string                 `json:"remark"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type ItemListResponse struct {
	Items    []ItemResponse `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

type DictTypeEntity = model.SystemDictType
type DictItemEntity = model.SystemDictItem

const (
	PermissionTypeList   = "system:dict:type:list"
	PermissionTypeCreate = "system:dict:type:create"
	PermissionTypeUpdate = "system:dict:type:update"
	PermissionTypeStatus = "system:dict:type:update_status"
	PermissionItemList   = "system:dict:item:list"
	PermissionItemCreate = "system:dict:item:create"
	PermissionItemUpdate = "system:dict:item:update"
	PermissionItemStatus = "system:dict:item:update_status"
)

func NormalizeTypeStatusFilter(value int) (*model.SystemDictStatus, error) {
	if value == 0 {
		return nil, nil
	}
	status := model.SystemDictStatus(value)
	if !ValidStatus(status) {
		return nil, errorsx.BadRequest("字典状态不正确")
	}
	return &status, nil
}

func NormalizeItemStatusFilter(value int) (*model.SystemDictStatus, error) {
	return NormalizeTypeStatusFilter(value)
}

func NormalizeCreateTypeRequest(req CreateTypeRequest) (CreateTypeRequest, error) {
	code, err := normalizeCode("字典编码", req.Code, 64)
	if err != nil {
		return CreateTypeRequest{}, err
	}
	name, err := normalizeName("字典名称", req.Name, 64)
	if err != nil {
		return CreateTypeRequest{}, err
	}
	status, err := NormalizeStatus(req.Status, true)
	if err != nil {
		return CreateTypeRequest{}, err
	}
	remark, err := normalizeRemark(req.Remark)
	if err != nil {
		return CreateTypeRequest{}, err
	}

	req.Code = code
	req.Name = name
	req.Status = status
	req.Remark = remark
	return req, nil
}

func NormalizeUpdateTypeRequest(req UpdateTypeRequest) (UpdateTypeRequest, error) {
	name, err := normalizeName("字典名称", req.Name, 64)
	if err != nil {
		return UpdateTypeRequest{}, err
	}
	status, err := NormalizeStatus(req.Status, false)
	if err != nil {
		return UpdateTypeRequest{}, err
	}
	remark, err := normalizeRemark(req.Remark)
	if err != nil {
		return UpdateTypeRequest{}, err
	}

	req.Name = name
	req.Status = status
	req.Remark = remark
	return req, nil
}

func NormalizeCreateItemRequest(req CreateItemRequest) (CreateItemRequest, error) {
	if req.TypeID == 0 {
		return CreateItemRequest{}, errorsx.BadRequest("字典类型 ID 不正确")
	}
	itemKey, err := normalizeCode("字典项编码", req.ItemKey, 64)
	if err != nil {
		return CreateItemRequest{}, err
	}
	label, err := normalizeName("字典项名称", req.Label, 64)
	if err != nil {
		return CreateItemRequest{}, err
	}
	value, err := normalizeRequiredText("字典项值", req.Value, 255)
	if err != nil {
		return CreateItemRequest{}, err
	}
	tagType, err := normalizeOptionalText("标签样式", req.TagType, 32)
	if err != nil {
		return CreateItemRequest{}, err
	}
	status, err := NormalizeStatus(req.Status, true)
	if err != nil {
		return CreateItemRequest{}, err
	}
	remark, err := normalizeRemark(req.Remark)
	if err != nil {
		return CreateItemRequest{}, err
	}

	req.ItemKey = itemKey
	req.Label = label
	req.Value = value
	req.TagType = tagType
	req.Status = status
	req.Remark = remark
	return req, nil
}

func NormalizeUpdateItemRequest(req UpdateItemRequest) (UpdateItemRequest, error) {
	label, err := normalizeName("字典项名称", req.Label, 64)
	if err != nil {
		return UpdateItemRequest{}, err
	}
	value, err := normalizeRequiredText("字典项值", req.Value, 255)
	if err != nil {
		return UpdateItemRequest{}, err
	}
	tagType, err := normalizeOptionalText("标签样式", req.TagType, 32)
	if err != nil {
		return UpdateItemRequest{}, err
	}
	status, err := NormalizeStatus(req.Status, false)
	if err != nil {
		return UpdateItemRequest{}, err
	}
	remark, err := normalizeRemark(req.Remark)
	if err != nil {
		return UpdateItemRequest{}, err
	}

	req.Label = label
	req.Value = value
	req.TagType = tagType
	req.Status = status
	req.Remark = remark
	return req, nil
}

func ValidStatus(status model.SystemDictStatus) bool {
	return status == model.SystemDictStatusEnabled || status == model.SystemDictStatusDisabled
}

func NormalizeStatus(status model.SystemDictStatus, allowDefault bool) (model.SystemDictStatus, error) {
	if status == 0 && allowDefault {
		status = model.SystemDictStatusEnabled
	}
	if !ValidStatus(status) {
		return 0, errorsx.BadRequest("字典状态不正确")
	}
	return status, nil
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

func normalizeName(fieldName string, value string, maxLen int) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errorsx.BadRequest(fieldName + "不能为空")
	}
	if len(value) > maxLen {
		return "", errorsx.BadRequest(fieldName + "长度不能超过 " + strconv.Itoa(maxLen) + " 个字符")
	}
	return value, nil
}

func normalizeRequiredText(fieldName string, value string, maxLen int) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errorsx.BadRequest(fieldName + "不能为空")
	}
	if len(value) > maxLen {
		return "", errorsx.BadRequest(fieldName + "长度不能超过 " + strconv.Itoa(maxLen) + " 个字符")
	}
	return value, nil
}

func normalizeOptionalText(fieldName string, value string, maxLen int) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) > maxLen {
		return "", errorsx.BadRequest(fieldName + "长度不能超过 " + strconv.Itoa(maxLen) + " 个字符")
	}
	return value, nil
}

func normalizeRemark(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) > 255 {
		return "", errorsx.BadRequest("备注不能超过 255 个字符")
	}
	return value, nil
}

func BuildTypeResponse(item DictTypeEntity) TypeResponse {
	return TypeResponse{
		ID:        item.ID,
		Code:      item.Code,
		Name:      item.Name,
		Sort:      item.Sort,
		Status:    item.Status,
		Remark:    item.Remark,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

func BuildItemResponse(item DictItemEntity) ItemResponse {
	return ItemResponse{
		ID:        item.ID,
		TypeID:    item.TypeID,
		ItemKey:   item.ItemKey,
		Label:     item.Label,
		Value:     item.Value,
		TagType:   item.TagType,
		Sort:      item.Sort,
		Status:    item.Status,
		Remark:    item.Remark,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}
