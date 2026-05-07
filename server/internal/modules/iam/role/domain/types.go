package domain

import (
	"strings"
	"time"

	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/datascope"
	"ez-admin-gin/server/internal/platform/model"
)

type ListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
	Status   int    `form:"status"`
}

type CreateRequest struct {
	Code                string           `json:"code"`
	Name                string           `json:"name"`
	Sort                int              `json:"sort"`
	DataScope           datascope.Scope  `json:"data_scope"`
	CustomDepartmentIDs []uint           `json:"custom_department_ids"`
	Status              model.RoleStatus `json:"status"`
	Remark              string           `json:"remark"`
}

type UpdateRequest struct {
	Name                string           `json:"name"`
	Sort                int              `json:"sort"`
	DataScope           datascope.Scope  `json:"data_scope"`
	CustomDepartmentIDs []uint           `json:"custom_department_ids"`
	Status              model.RoleStatus `json:"status"`
	Remark              string           `json:"remark"`
}

type UpdateStatusRequest struct {
	Status model.RoleStatus `json:"status"`
}

type PermissionItem struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

type UpdatePermissionsRequest struct {
	Permissions []PermissionItem `json:"permissions"`
}

type UpdateMenusRequest struct {
	MenuIDs []uint `json:"menu_ids"`
}

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
	MenuIDs             []uint           `json:"menu_ids"`
	CreatedAt           time.Time        `json:"created_at"`
	UpdatedAt           time.Time        `json:"updated_at"`
}

type ListResponse struct {
	Items    []Response `json:"items"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

type Entity = model.Role
type RoleMenuEntity = model.RoleMenu
type RoleDataScopeEntity = model.RoleDataScope

const (
	PermissionList              = "system:role:list"
	PermissionCreate            = "system:role:create"
	PermissionUpdate            = "system:role:update"
	PermissionUpdateStatus      = "system:role:update_status"
	PermissionUpdatePermissions = "system:role:update_permissions"
	PermissionUpdateMenus       = "system:role:update_menus"
	SuperAdminRoleCode          = "super_admin"
)

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

func NormalizePermissions(permissions []PermissionItem) ([]PermissionItem, error) {
	unique := make([]PermissionItem, 0, len(permissions))
	seen := make(map[string]struct{}, len(permissions))

	for _, item := range permissions {
		path := strings.TrimSpace(item.Path)
		method := strings.ToUpper(strings.TrimSpace(item.Method))
		if path == "" || method == "" {
			return nil, errorsx.BadRequest("接口权限参数不正确")
		}

		key := path + " " + method
		if _, ok := seen[key]; ok {
			continue
		}

		seen[key] = struct{}{}
		unique = append(unique, PermissionItem{Path: path, Method: method})
	}

	return unique, nil
}

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

func ValidRoleStatus(status model.RoleStatus) bool {
	return status == model.RoleStatusEnabled || status == model.RoleStatusDisabled
}

func ValidDataScope(scope datascope.Scope) bool {
	switch scope {
	case datascope.ScopeAll, datascope.ScopeDept, datascope.ScopeDeptAndChildren, datascope.ScopeSelf, datascope.ScopeCustomDept:
		return true
	default:
		return false
	}
}

func BuildResponse(role model.Role, customDepartmentIDs []uint, permissions []PermissionItem, menuIDs []uint) Response {
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
		MenuIDs:             menuIDs,
		CreatedAt:           role.CreatedAt,
		UpdatedAt:           role.UpdatedAt,
	}
}
