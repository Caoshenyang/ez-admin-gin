// Package domain 定义部门的请求/响应结构、权限常量和业务校验规则。
package domain

import (
	"fmt"
	"strings"
	"time"

	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/model"
)

// ListQuery 定义部门列表查询的过滤参数。
type ListQuery struct {
	Keyword string `form:"keyword"`
	Status  int    `form:"status"`
}

// CreateRequest 定义创建部门的请求参数。
type CreateRequest struct {
	ParentID     uint                   `json:"parent_id"`
	Name         string                 `json:"name"`
	Code         string                 `json:"code"`
	LeaderUserID uint                   `json:"leader_user_id"`
	Sort         int                    `json:"sort"`
	Status       model.DepartmentStatus `json:"status"`
	Remark       string                 `json:"remark"`
}

// UpdateRequest 定义更新部门的请求参数。
type UpdateRequest struct {
	ParentID     uint                   `json:"parent_id"`
	Name         string                 `json:"name"`
	Code         string                 `json:"code"`
	LeaderUserID uint                   `json:"leader_user_id"`
	Sort         int                    `json:"sort"`
	Status       model.DepartmentStatus `json:"status"`
	Remark       string                 `json:"remark"`
}

// UpdateStatusRequest 定义切换部门状态的请求参数。
type UpdateStatusRequest struct {
	Status model.DepartmentStatus `json:"status"`
}

// Response 定义部门信息的响应结构。
type Response struct {
	ID           uint                   `json:"id"`
	ParentID     uint                   `json:"parent_id"`
	Ancestors    string                 `json:"ancestors"`
	Name         string                 `json:"name"`
	Code         string                 `json:"code"`
	LeaderUserID uint                   `json:"leader_user_id"`
	Sort         int                    `json:"sort"`
	Status       model.DepartmentStatus `json:"status"`
	Remark       string                 `json:"remark"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	Children     []Response             `json:"children,omitempty"`
}

// Entity 是部门聚合根模型的类型别名。
type Entity = model.Department

const (
	PermissionList         = "system:department:list"
	PermissionCreate       = "system:department:create"
	PermissionUpdate       = "system:department:update"
	PermissionUpdateStatus = "system:department:status"
)

// NormalizeDepartmentInput 规范化并校验部门输入参数。
func NormalizeDepartmentInput(parentID uint, name string, code string, leaderUserID uint, sort int, status model.DepartmentStatus, remark string) (uint, string, string, uint, int, model.DepartmentStatus, string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return 0, "", "", 0, 0, 0, "", errorsx.BadRequest("部门名称不能为空")
	}
	if len(name) > 64 {
		return 0, "", "", 0, 0, 0, "", errorsx.BadRequest("部门名称不能超过 64 个字符")
	}

	code = strings.TrimSpace(code)
	if code == "" {
		return 0, "", "", 0, 0, 0, "", errorsx.BadRequest("部门编码不能为空")
	}
	if len(code) > 64 {
		return 0, "", "", 0, 0, 0, "", errorsx.BadRequest("部门编码不能超过 64 个字符")
	}

	if status == 0 {
		status = model.DepartmentStatusEnabled
	}
	if !ValidStatus(status) {
		return 0, "", "", 0, 0, 0, "", errorsx.BadRequest("部门状态不正确")
	}

	remark = strings.TrimSpace(remark)
	if len(remark) > 255 {
		return 0, "", "", 0, 0, 0, "", errorsx.BadRequest("备注不能超过 255 个字符")
	}

	return parentID, name, code, leaderUserID, sort, status, remark, nil
}

// ValidStatus 判断部门状态值是否合法。
func ValidStatus(status model.DepartmentStatus) bool {
	return status == model.DepartmentStatusEnabled || status == model.DepartmentStatusDisabled
}

// BuildResponse 将部门模型转换为响应结构。
func BuildResponse(item model.Department) Response {
	return Response{
		ID:           item.ID,
		ParentID:     item.ParentID,
		Ancestors:    item.Ancestors,
		Name:         item.Name,
		Code:         item.Code,
		LeaderUserID: item.LeaderUserID,
		Sort:         item.Sort,
		Status:       item.Status,
		Remark:       item.Remark,
		CreatedAt:    item.CreatedAt,
		UpdatedAt:    item.UpdatedAt,
	}
}

// BuildAncestors 根据父部门生成祖先路径字符串。
func BuildAncestors(parent model.Department) string {
	if parent.ID == 0 {
		return "0"
	}
	return fmt.Sprintf("%s,%d", parent.Ancestors, parent.ID)
}

// FullPath 返回包含自身 ID 的完整祖先路径。
func FullPath(item model.Department) string {
	if item.Ancestors == "" {
		return fmt.Sprintf("%d", item.ID)
	}
	return fmt.Sprintf("%s,%d", item.Ancestors, item.ID)
}

// IsDescendantPath 判断 path 是否在 target 子树内。
func IsDescendantPath(path string, target string) bool {
	return path == target || strings.HasPrefix(path, target+",")
}
