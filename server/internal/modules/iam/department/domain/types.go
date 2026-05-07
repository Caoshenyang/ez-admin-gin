package domain

import (
	"strings"
	"time"

	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/model"
)

type ListQuery struct {
	Keyword string `form:"keyword"`
	Status  int    `form:"status"`
}

type CreateRequest struct {
	ParentID     uint                   `json:"parent_id"`
	Name         string                 `json:"name"`
	Code         string                 `json:"code"`
	LeaderUserID uint                   `json:"leader_user_id"`
	Sort         int                    `json:"sort"`
	Status       model.DepartmentStatus `json:"status"`
	Remark       string                 `json:"remark"`
}

type UpdateRequest struct {
	ParentID     uint                   `json:"parent_id"`
	Name         string                 `json:"name"`
	Code         string                 `json:"code"`
	LeaderUserID uint                   `json:"leader_user_id"`
	Sort         int                    `json:"sort"`
	Status       model.DepartmentStatus `json:"status"`
	Remark       string                 `json:"remark"`
}

type UpdateStatusRequest struct {
	Status model.DepartmentStatus `json:"status"`
}

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

type Entity = model.Department

const (
	PermissionList         = "system:department:list"
	PermissionCreate       = "system:department:create"
	PermissionUpdate       = "system:department:update"
	PermissionUpdateStatus = "system:department:status"
)

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

func ValidStatus(status model.DepartmentStatus) bool {
	return status == model.DepartmentStatusEnabled || status == model.DepartmentStatusDisabled
}

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
