package customer

import (
	"strconv"
	"strings"
	"time"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
)

// ListQuery 表示客户分页查询参数。
type ListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
	Level    string `form:"level"`
	Source   string `form:"source"`
	Status   int    `form:"status"`
}

// CreateRequest 表示创建客户请求体。
type CreateRequest struct {
	Name        string               `json:"name"`
	ContactName string               `json:"contact_name"`
	Phone       string               `json:"phone"`
	Level       string               `json:"level"`
	Source      string               `json:"source"`
	Status      model.CustomerStatus `json:"status"`
	Remark      string               `json:"remark"`
}

// UpdateRequest 表示编辑客户请求体。
type UpdateRequest struct {
	Name        string               `json:"name"`
	ContactName string               `json:"contact_name"`
	Phone       string               `json:"phone"`
	Level       string               `json:"level"`
	Source      string               `json:"source"`
	Status      model.CustomerStatus `json:"status"`
	Remark      string               `json:"remark"`
}

// UpdateStatusRequest 表示单独修改客户状态请求体。
type UpdateStatusRequest struct {
	Status model.CustomerStatus `json:"status"`
}

// Response 表示客户对象返回结构。
type Response struct {
	ID             uint                 `json:"id"`
	Name           string               `json:"name"`
	ContactName    string               `json:"contact_name"`
	Phone          string               `json:"phone"`
	Level          string               `json:"level"`
	Source         string               `json:"source"`
	DepartmentID   uint                 `json:"department_id"`
	DepartmentName string               `json:"department_name"`
	OwnerUserID    uint                 `json:"owner_user_id"`
	OwnerUsername  string               `json:"owner_username"`
	OwnerNickname  string               `json:"owner_nickname"`
	Status         model.CustomerStatus `json:"status"`
	Remark         string               `json:"remark"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
}

// ListResponse 表示客户分页结果。
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
	name, err := normalizeName(req.Name)
	if err != nil {
		return CreateRequest{}, err
	}
	status, err := normalizeStatus(req.Status, true)
	if err != nil {
		return CreateRequest{}, err
	}
	contactName, err := normalizeContactName(req.ContactName)
	if err != nil {
		return CreateRequest{}, err
	}
	phone, err := normalizePhone(req.Phone)
	if err != nil {
		return CreateRequest{}, err
	}
	level, err := normalizeLevel(req.Level)
	if err != nil {
		return CreateRequest{}, err
	}
	source, err := normalizeSource(req.Source)
	if err != nil {
		return CreateRequest{}, err
	}
	remark, err := normalizeRemark(req.Remark)
	if err != nil {
		return CreateRequest{}, err
	}

	req.Name = name
	req.ContactName = contactName
	req.Phone = phone
	req.Level = level
	req.Source = source
	req.Status = status
	req.Remark = remark
	return req, nil
}

// NormalizeUpdateRequest 统一校验并收敛编辑参数。
func NormalizeUpdateRequest(req UpdateRequest) (UpdateRequest, error) {
	name, err := normalizeName(req.Name)
	if err != nil {
		return UpdateRequest{}, err
	}
	status, err := normalizeStatus(req.Status, false)
	if err != nil {
		return UpdateRequest{}, err
	}
	contactName, err := normalizeContactName(req.ContactName)
	if err != nil {
		return UpdateRequest{}, err
	}
	phone, err := normalizePhone(req.Phone)
	if err != nil {
		return UpdateRequest{}, err
	}
	level, err := normalizeLevel(req.Level)
	if err != nil {
		return UpdateRequest{}, err
	}
	source, err := normalizeSource(req.Source)
	if err != nil {
		return UpdateRequest{}, err
	}
	remark, err := normalizeRemark(req.Remark)
	if err != nil {
		return UpdateRequest{}, err
	}

	req.Name = name
	req.ContactName = contactName
	req.Phone = phone
	req.Level = level
	req.Source = source
	req.Status = status
	req.Remark = remark
	return req, nil
}

// NormalizeStatusFilter 把状态查询参数转换成客户状态。
func NormalizeStatusFilter(value int) (*model.CustomerStatus, error) {
	if value == 0 {
		return nil, nil
	}

	status := model.CustomerStatus(value)
	if !validStatus(status) {
		return nil, apperror.BadRequest("客户状态不正确")
	}

	return &status, nil
}

// ParseCustomerID 解析路径参数中的客户 ID。
func ParseCustomerID(value string) (uint, error) {
	id, err := strconv.ParseUint(strings.TrimSpace(value), 10, 64)
	if err != nil || id == 0 {
		return 0, apperror.BadRequest("客户 ID 不正确")
	}
	return uint(id), nil
}

// BuildResponse 把联合查询视图压成 API 返回结构。
func BuildResponse(item View) Response {
	return Response{
		ID:             item.ID,
		Name:           item.Name,
		ContactName:    item.ContactName,
		Phone:          item.Phone,
		Level:          item.Level,
		Source:         item.Source,
		DepartmentID:   item.DepartmentID,
		DepartmentName: item.DepartmentName,
		OwnerUserID:    item.OwnerUserID,
		OwnerUsername:  item.OwnerUsername,
		OwnerNickname:  item.OwnerNickname,
		Status:         item.Status,
		Remark:         item.Remark,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
	}
}

func normalizeName(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", apperror.BadRequest("客户名称不能为空")
	}
	if len(value) > 128 {
		return "", apperror.BadRequest("客户名称不能超过 128 个字符")
	}
	return value, nil
}

func normalizeContactName(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) > 64 {
		return "", apperror.BadRequest("联系人不能超过 64 个字符")
	}
	return value, nil
}

func normalizePhone(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) > 32 {
		return "", apperror.BadRequest("联系电话不能超过 32 个字符")
	}
	return value, nil
}

func normalizeLevel(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) > 32 {
		return "", apperror.BadRequest("客户等级不能超过 32 个字符")
	}
	return value, nil
}

func normalizeSource(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) > 32 {
		return "", apperror.BadRequest("客户来源不能超过 32 个字符")
	}
	return value, nil
}

func normalizeRemark(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) > 255 {
		return "", apperror.BadRequest("备注不能超过 255 个字符")
	}
	return value, nil
}

func normalizeStatus(status model.CustomerStatus, allowDefault bool) (model.CustomerStatus, error) {
	if status == 0 && allowDefault {
		status = model.CustomerStatusEnabled
	}
	if !validStatus(status) {
		return 0, apperror.BadRequest("客户状态不正确")
	}
	return status, nil
}

func validStatus(status model.CustomerStatus) bool {
	return status == model.CustomerStatusEnabled || status == model.CustomerStatusDisabled
}
