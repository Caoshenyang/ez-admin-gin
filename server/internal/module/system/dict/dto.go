package dict

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
)

var codePattern = regexp.MustCompile(`^[a-z0-9:_-]+$`)

// TypeListQuery 表示字典类型分页查询参数。
type TypeListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
	Status   int    `form:"status"`
}

// CreateTypeRequest 表示创建字典类型请求体。
type CreateTypeRequest struct {
	Code   string                 `json:"code"`
	Name   string                 `json:"name"`
	Sort   int                    `json:"sort"`
	Status model.SystemDictStatus `json:"status"`
	Remark string                 `json:"remark"`
}

// UpdateTypeRequest 表示编辑字典类型请求体。
type UpdateTypeRequest struct {
	Name   string                 `json:"name"`
	Sort   int                    `json:"sort"`
	Status model.SystemDictStatus `json:"status"`
	Remark string                 `json:"remark"`
}

// UpdateTypeStatusRequest 表示字典类型状态更新请求体。
type UpdateTypeStatusRequest struct {
	Status model.SystemDictStatus `json:"status"`
}

// TypeResponse 表示字典类型返回结构。
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

// TypeListResponse 表示字典类型分页结果。
type TypeListResponse struct {
	Items    []TypeResponse `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

// ItemListQuery 表示字典项分页查询参数。
type ItemListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	TypeID   uint   `form:"type_id"`
	Keyword  string `form:"keyword"`
	Status   int    `form:"status"`
}

// CreateItemRequest 表示创建字典项请求体。
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

// UpdateItemRequest 表示编辑字典项请求体。
type UpdateItemRequest struct {
	Label   string                 `json:"label"`
	Value   string                 `json:"value"`
	TagType string                 `json:"tag_type"`
	Sort    int                    `json:"sort"`
	Status  model.SystemDictStatus `json:"status"`
	Remark  string                 `json:"remark"`
}

// UpdateItemStatusRequest 表示字典项状态更新请求体。
type UpdateItemStatusRequest struct {
	Status model.SystemDictStatus `json:"status"`
}

// ItemResponse 表示字典项返回结构。
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

// ItemListResponse 表示字典项分页结果。
type ItemListResponse struct {
	Items    []ItemResponse `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
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

// NormalizeTypeStatusFilter 把状态查询参数转换成字典状态。
func NormalizeTypeStatusFilter(value int) (*model.SystemDictStatus, error) {
	if value == 0 {
		return nil, nil
	}

	status := model.SystemDictStatus(value)
	if !ValidStatus(status) {
		return nil, apperror.BadRequest("字典状态不正确")
	}

	return &status, nil
}

// NormalizeItemStatusFilter 把状态查询参数转换成字典状态。
func NormalizeItemStatusFilter(value int) (*model.SystemDictStatus, error) {
	return NormalizeTypeStatusFilter(value)
}

// NormalizeCreateTypeRequest 统一校验并收敛创建字典类型参数。
func NormalizeCreateTypeRequest(req CreateTypeRequest) (CreateTypeRequest, error) {
	code, err := normalizeCode("字典编码", req.Code, 64)
	if err != nil {
		return CreateTypeRequest{}, err
	}
	name, err := normalizeName("字典名称", req.Name, 64)
	if err != nil {
		return CreateTypeRequest{}, err
	}
	status, err := normalizeStatus(req.Status, true)
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

// NormalizeUpdateTypeRequest 统一校验并收敛编辑字典类型参数。
func NormalizeUpdateTypeRequest(req UpdateTypeRequest) (UpdateTypeRequest, error) {
	name, err := normalizeName("字典名称", req.Name, 64)
	if err != nil {
		return UpdateTypeRequest{}, err
	}
	status, err := normalizeStatus(req.Status, false)
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

// NormalizeCreateItemRequest 统一校验并收敛创建字典项参数。
func NormalizeCreateItemRequest(req CreateItemRequest) (CreateItemRequest, error) {
	if req.TypeID == 0 {
		return CreateItemRequest{}, apperror.BadRequest("字典类型 ID 不正确")
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
	status, err := normalizeStatus(req.Status, true)
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

// NormalizeUpdateItemRequest 统一校验并收敛编辑字典项参数。
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
	status, err := normalizeStatus(req.Status, false)
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

// ValidStatus 判断字典状态是否合法。
func ValidStatus(status model.SystemDictStatus) bool {
	return status == model.SystemDictStatusEnabled || status == model.SystemDictStatusDisabled
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

func normalizeName(fieldName string, value string, maxLen int) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", apperror.BadRequest(fieldName + "不能为空")
	}
	if len(value) > maxLen {
		return "", apperror.BadRequest(fieldName + "长度不能超过 " + strconv.Itoa(maxLen) + " 个字符")
	}

	return value, nil
}

func normalizeRequiredText(fieldName string, value string, maxLen int) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", apperror.BadRequest(fieldName + "不能为空")
	}
	if len(value) > maxLen {
		return "", apperror.BadRequest(fieldName + "长度不能超过 " + strconv.Itoa(maxLen) + " 个字符")
	}

	return value, nil
}

func normalizeOptionalText(fieldName string, value string, maxLen int) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) > maxLen {
		return "", apperror.BadRequest(fieldName + "长度不能超过 " + strconv.Itoa(maxLen) + " 个字符")
	}

	return value, nil
}

func normalizeStatus(status model.SystemDictStatus, allowDefault bool) (model.SystemDictStatus, error) {
	if status == 0 && allowDefault {
		status = model.SystemDictStatusEnabled
	}
	if !ValidStatus(status) {
		return 0, apperror.BadRequest("字典状态不正确")
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

// BuildTypeResponse 把字典类型模型压成 API 返回结构。
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

// BuildItemResponse 把字典项模型压成 API 返回结构。
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
