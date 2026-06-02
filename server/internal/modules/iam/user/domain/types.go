// Package domain 定义用户的请求/响应结构、权限常量和业务校验规则。
package domain

import (
	"strings"
	"time"

	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/model"
)

// ListQuery 定义用户列表查询的过滤参数。
type ListQuery struct {
	Page         int    `form:"page"`
	PageSize     int    `form:"page_size"`
	Keyword      string `form:"keyword"`
	Status       int    `form:"status"`
	DepartmentID uint   `form:"department_id"`
	RoleID       uint   `form:"role_id"`
}

// CreateRequest 定义创建用户的请求参数。
type CreateRequest struct {
	Username     string           `json:"username"`
	Password     string           `json:"password"`
	Nickname     string           `json:"nickname"`
	DepartmentID uint             `json:"department_id"`
	Status       model.UserStatus `json:"status"`
	RoleIDs      []uint           `json:"role_ids"`
	PostIDs      []uint           `json:"post_ids"`
}

// UpdateRequest 定义更新用户的请求参数。
type UpdateRequest struct {
	Nickname     string           `json:"nickname"`
	DepartmentID uint             `json:"department_id"`
	Status       model.UserStatus `json:"status"`
	PostIDs      []uint           `json:"post_ids"`
}

// UpdateStatusRequest 定义切换用户状态的请求参数。
type UpdateStatusRequest struct {
	Status model.UserStatus `json:"status"`
}

// UpdateRolesRequest 定义更新用户角色的请求参数。
type UpdateRolesRequest struct {
	RoleIDs []uint `json:"role_ids"`
}

// Response 定义用户信息的响应结构。
type Response struct {
	ID           uint             `json:"id"`
	Username     string           `json:"username"`
	Nickname     string           `json:"nickname"`
	DepartmentID uint             `json:"department_id"`
	Status       model.UserStatus `json:"status"`
	RoleIDs      []uint           `json:"role_ids"`
	PostIDs      []uint           `json:"post_ids"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

// ListResponse 定义用户分页列表的响应结构。
type ListResponse struct {
	Items    []Response `json:"items"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

// Entity 是用户聚合根模型的类型别名。
type Entity = model.User

// UserRoleEntity 是用户角色关联模型的类型别名。
type UserRoleEntity = model.UserRole

const (
	PermissionList         = "system:user:list"
	PermissionCreate       = "system:user:create"
	PermissionUpdate       = "system:user:update"
	PermissionUpdateStatus = "system:user:update_status"
	PermissionUpdateRoles  = "system:user:update_roles"
)

// NormalizeCreateRequest 规范化并校验用户创建请求参数。
func NormalizeCreateRequest(req CreateRequest) (CreateRequest, error) {
	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" {
		return CreateRequest{}, errorsx.BadRequest("用户名不能为空")
	}
	if len(req.Username) > 64 {
		return CreateRequest{}, errorsx.BadRequest("用户名不能超过 64 个字符")
	}

	if len(req.Password) < 8 || len(req.Password) > 72 {
		return CreateRequest{}, errorsx.BadRequest("密码长度需要在 8 到 72 个字符之间")
	}

	req.Nickname = strings.TrimSpace(req.Nickname)
	if req.Nickname == "" {
		req.Nickname = req.Username
	}
	if len(req.Nickname) > 64 {
		return CreateRequest{}, errorsx.BadRequest("昵称不能超过 64 个字符")
	}

	if req.Status == 0 {
		req.Status = model.UserStatusEnabled
	}
	if !ValidStatus(req.Status) {
		return CreateRequest{}, errorsx.BadRequest("用户状态不正确")
	}

	roleIDs, err := NormalizeRoleIDs(req.RoleIDs)
	if err != nil {
		return CreateRequest{}, err
	}
	req.RoleIDs = roleIDs
	postIDs, err := NormalizePostIDs(req.PostIDs)
	if err != nil {
		return CreateRequest{}, err
	}
	req.PostIDs = postIDs

	return req, nil
}

// NormalizeUpdateRequest 规范化并校验用户更新请求参数。
func NormalizeUpdateRequest(req UpdateRequest) (UpdateRequest, error) {
	req.Nickname = strings.TrimSpace(req.Nickname)
	if req.Nickname == "" {
		return UpdateRequest{}, errorsx.BadRequest("昵称不能为空")
	}
	if len(req.Nickname) > 64 {
		return UpdateRequest{}, errorsx.BadRequest("昵称不能超过 64 个字符")
	}
	if !ValidStatus(req.Status) {
		return UpdateRequest{}, errorsx.BadRequest("用户状态不正确")
	}
	postIDs, err := NormalizePostIDs(req.PostIDs)
	if err != nil {
		return UpdateRequest{}, err
	}
	req.PostIDs = postIDs

	return req, nil
}

// NormalizeRoleIDs 去重并校验角色 ID 列表。
func NormalizeRoleIDs(roleIDs []uint) ([]uint, error) {
	return normalizeUintIDs(roleIDs, "角色 ID 不正确")
}

// NormalizePostIDs 去重并校验岗位 ID 列表。
func NormalizePostIDs(postIDs []uint) ([]uint, error) {
	return normalizeUintIDs(postIDs, "岗位 ID 不正确")
}

func normalizeUintIDs(ids []uint, invalidMessage string) ([]uint, error) {
	unique := make([]uint, 0, len(ids))
	seen := make(map[uint]struct{}, len(ids))

	for _, id := range ids {
		if id == 0 {
			return nil, errorsx.BadRequest(invalidMessage)
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		unique = append(unique, id)
	}

	return unique, nil
}

// ValidStatus 判断用户状态值是否合法。
func ValidStatus(status model.UserStatus) bool {
	return status == model.UserStatusEnabled || status == model.UserStatusDisabled
}

// BuildResponse 将用户模型及关联数据转换为响应结构。
func BuildResponse(user model.User, roleIDs []uint, postIDs []uint) Response {
	return Response{
		ID:           user.ID,
		Username:     user.Username,
		Nickname:     user.Nickname,
		DepartmentID: user.DepartmentID,
		Status:       user.Status,
		RoleIDs:      roleIDs,
		PostIDs:      postIDs,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}
