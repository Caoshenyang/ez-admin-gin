// Package domain 定义角色的请求/响应结构、权限常量和业务校验规则。
package domain

import (
	"strings"
	"time"

	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/datascope"
	"ez-admin-gin/server/internal/platform/model"
)

// ListQuery 定义角色列表查询的过滤参数。
type ListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
	Status   int    `form:"status"`
}

// CreateRequest 定义创建角色的请求参数。
type CreateRequest struct {
	Code                string           `json:"code"`
	Name                string           `json:"name"`
	Sort                int              `json:"sort"`
	DataScope           datascope.Scope  `json:"data_scope"`
	CustomDepartmentIDs []uint           `json:"custom_department_ids"`
	Status              model.RoleStatus `json:"status"`
	Remark              string           `json:"remark"`
}

// UpdateRequest 定义更新角色的请求参数。
type UpdateRequest struct {
	Name                string           `json:"name"`
	Sort                int              `json:"sort"`
	DataScope           datascope.Scope  `json:"data_scope"`
	CustomDepartmentIDs []uint           `json:"custom_department_ids"`
	Status              model.RoleStatus `json:"status"`
	Remark              string           `json:"remark"`
}

// UpdateStatusRequest 定义切换角色状态的请求参数。
type UpdateStatusRequest struct {
	Status model.RoleStatus `json:"status"`
}

// PermissionItem 表示一条接口权限元数据快照。
type PermissionItem struct {
	ID     uint            `json:"id"`
	Code   string          `json:"code"`
	Name   string          `json:"name"`
	Module string          `json:"module"`
	Method string          `json:"method"`
	Path   string          `json:"path"`
	Status model.APIStatus `json:"status"`
}

// UpdatePermissionsRequest 定义更新角色权限的请求参数。
type UpdatePermissionsRequest struct {
	APIIDs []uint `json:"api_ids"`
}

// UpdateMenusRequest 定义更新角色菜单的请求参数。
type UpdateMenusRequest struct {
	MenuIDs []uint `json:"menu_ids"`
}

// Response 定义角色信息的响应结构。
type Response struct {
	ID                  uint             `json:"id"`
	Code                string           `json:"code"`
	Name                string           `json:"name"`
	Sort                int              `json:"sort"`
	DataScope           datascope.Scope  `json:"data_scope"`
	CustomDepartmentIDs []uint           `json:"custom_department_ids"`
	Status              model.RoleStatus `json:"status"`
	Remark              string           `json:"remark"`
	Permissions         []PermissionItem `json:"permissions"`
	APIIDs              []uint           `json:"api_ids"`
	MenuIDs             []uint           `json:"menu_ids"`
	CreatedAt           time.Time        `json:"created_at"`
	UpdatedAt           time.Time        `json:"updated_at"`
}

// ListResponse 定义角色分页列表的响应结构。
type ListResponse struct {
	Items    []Response `json:"items"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

// Entity 是角色聚合根模型的类型别名。
type Entity = model.Role

// RoleMenuEntity 是角色菜单关联模型的类型别名。
type RoleMenuEntity = model.RoleMenu

// RoleDataScopeEntity 是角色数据范围关联模型的类型别名。
type RoleDataScopeEntity = model.RoleDataScope

const (
	PermissionList              = "system:role:list"
	PermissionCreate            = "system:role:create"
	PermissionUpdate            = "system:role:update"
	PermissionUpdateStatus      = "system:role:status"
	PermissionUpdatePermissions = "system:role:permission"
	PermissionUpdateMenus       = "system:role:menu"
	PermissionDelete            = "system:role:delete"
	SuperAdminRoleCode          = "super_admin"
)

// NormalizeCreateRequest 规范化并校验角色创建请求参数。
func NormalizeCreateRequest(req CreateRequest) (CreateRequest, error) {
	req.Code = strings.TrimSpace(req.Code)
	if req.Code == "" {
		return CreateRequest{}, errorsx.BadRequest("角色编码不能为空")
	}
	if len(req.Code) > 64 {
		return CreateRequest{}, errorsx.BadRequest("角色编码不能超过 64 个字符")
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return CreateRequest{}, errorsx.BadRequest("角色名称不能为空")
	}
	if len(req.Name) > 64 {
		return CreateRequest{}, errorsx.BadRequest("角色名称不能超过 64 个字符")
	}

	if req.DataScope == "" {
		req.DataScope = datascope.ScopeSelf
	}
	if !ValidDataScope(req.DataScope) {
		return CreateRequest{}, errorsx.BadRequest("角色数据范围不正确")
	}

	if req.Status == 0 {
		req.Status = model.RoleStatusEnabled
	}
	if !ValidRoleStatus(req.Status) {
		return CreateRequest{}, errorsx.BadRequest("角色状态不正确")
	}

	req.Remark = strings.TrimSpace(req.Remark)
	if len(req.Remark) > 255 {
		return CreateRequest{}, errorsx.BadRequest("备注不能超过 255 个字符")
	}

	customDepartmentIDs, err := NormalizeIDs(req.CustomDepartmentIDs, "部门 ID 不正确")
	if err != nil {
		return CreateRequest{}, err
	}
	if req.DataScope != datascope.ScopeCustomDept {
		customDepartmentIDs = nil
	}
	req.CustomDepartmentIDs = customDepartmentIDs

	return req, nil
}

// NormalizeUpdateRequest 规范化并校验角色更新请求参数。
func NormalizeUpdateRequest(req UpdateRequest) (UpdateRequest, error) {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return UpdateRequest{}, errorsx.BadRequest("角色名称不能为空")
	}
	if len(req.Name) > 64 {
		return UpdateRequest{}, errorsx.BadRequest("角色名称不能超过 64 个字符")
	}
	if !ValidRoleStatus(req.Status) {
		return UpdateRequest{}, errorsx.BadRequest("角色状态不正确")
	}
	if !ValidDataScope(req.DataScope) {
		return UpdateRequest{}, errorsx.BadRequest("角色数据范围不正确")
	}

	req.Remark = strings.TrimSpace(req.Remark)
	if len(req.Remark) > 255 {
		return UpdateRequest{}, errorsx.BadRequest("备注不能超过 255 个字符")
	}

	customDepartmentIDs, err := NormalizeIDs(req.CustomDepartmentIDs, "部门 ID 不正确")
	if err != nil {
		return UpdateRequest{}, err
	}
	if req.DataScope != datascope.ScopeCustomDept {
		customDepartmentIDs = nil
	}
	req.CustomDepartmentIDs = customDepartmentIDs

	return req, nil
}

// NormalizeIDs 去重并校验 ID 列表。
func NormalizeIDs(ids []uint, badRequestMessage string) ([]uint, error) {
	unique := make([]uint, 0, len(ids))
	seen := make(map[uint]struct{}, len(ids))

	for _, id := range ids {
		if id == 0 {
			return nil, errorsx.BadRequest(badRequestMessage)
		}
		if _, ok := seen[id]; ok {
			continue
		}

		seen[id] = struct{}{}
		unique = append(unique, id)
	}

	return unique, nil
}

// ValidRoleStatus 判断角色状态值是否合法。
func ValidRoleStatus(status model.RoleStatus) bool {
	return status == model.RoleStatusEnabled || status == model.RoleStatusDisabled
}

// ValidDataScope 判断数据范围值是否合法。
func ValidDataScope(scope datascope.Scope) bool {
	switch scope {
	case datascope.ScopeAll, datascope.ScopeDept, datascope.ScopeDeptAndChildren, datascope.ScopeSelf, datascope.ScopeCustomDept:
		return true
	default:
		return false
	}
}

// BuildResponse 将角色模型及关联数据转换为响应结构。
func BuildResponse(role model.Role, customDepartmentIDs []uint, permissions []PermissionItem, apiIDs []uint, menuIDs []uint) Response {
	return Response{
		ID:                  role.ID,
		Code:                role.Code,
		Name:                role.Name,
		Sort:                role.Sort,
		DataScope:           role.DataScope,
		CustomDepartmentIDs: customDepartmentIDs,
		Status:              role.Status,
		Remark:              role.Remark,
		Permissions:         permissions,
		APIIDs:              apiIDs,
		MenuIDs:             menuIDs,
		CreatedAt:           role.CreatedAt,
		UpdatedAt:           role.UpdatedAt,
	}
}
