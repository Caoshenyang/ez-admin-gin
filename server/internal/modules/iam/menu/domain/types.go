package domain

import (
	"strings"
	"time"

	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/model"
)

type CreateRequest struct {
	ParentID  uint             `json:"parent_id"`
	Type      model.MenuType   `json:"type"`
	Code      string           `json:"code"`
	Title     string           `json:"title"`
	Path      string           `json:"path"`
	Component string           `json:"component"`
	Icon      string           `json:"icon"`
	Sort      int              `json:"sort"`
	Status    model.MenuStatus `json:"status"`
	Remark    string           `json:"remark"`
}

type UpdateRequest struct {
	ParentID  uint             `json:"parent_id"`
	Type      model.MenuType   `json:"type"`
	Title     string           `json:"title"`
	Path      string           `json:"path"`
	Component string           `json:"component"`
	Icon      string           `json:"icon"`
	Sort      int              `json:"sort"`
	Status    model.MenuStatus `json:"status"`
	Remark    string           `json:"remark"`
}

type UpdateStatusRequest struct {
	Status model.MenuStatus `json:"status"`
}

type Response struct {
	ID        uint             `json:"id"`
	ParentID  uint             `json:"parent_id"`
	Type      model.MenuType   `json:"type"`
	Code      string           `json:"code"`
	Title     string           `json:"title"`
	Path      string           `json:"path"`
	Component string           `json:"component"`
	Icon      string           `json:"icon"`
	Sort      int              `json:"sort"`
	Status    model.MenuStatus `json:"status"`
	Remark    string           `json:"remark"`
	Children  []Response       `json:"children,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type Entity = model.Menu
type RoleMenuEntity = model.RoleMenu

const (
	PermissionList         = "system:menu:list"
	PermissionCreate       = "system:menu:create"
	PermissionUpdate       = "system:menu:update"
	PermissionUpdateStatus = "system:menu:update_status"
	PermissionDelete       = "system:menu:delete"
)

func NormalizeCreateRequest(req CreateRequest) (CreateRequest, error) {
	code := strings.TrimSpace(req.Code)
	if code == "" {
		return CreateRequest{}, errorsx.BadRequest("菜单编码不能为空")
	}
	if len(code) > 128 {
		return CreateRequest{}, errorsx.BadRequest("菜单编码不能超过 128 个字符")
	}

	title, path, component, icon, status, remark, err := normalizeFields(
		req.Type, req.Title, req.Path, req.Component, req.Icon, req.Status, req.Remark,
	)
	if err != nil {
		return CreateRequest{}, err
	}

	req.Code = code
	req.Title = title
	req.Path = path
	req.Component = component
	req.Icon = icon
	req.Status = status
	req.Remark = remark
	return req, nil
}

func NormalizeUpdateRequest(req UpdateRequest) (UpdateRequest, error) {
	title, path, component, icon, status, remark, err := normalizeFields(
		req.Type, req.Title, req.Path, req.Component, req.Icon, req.Status, req.Remark,
	)
	if err != nil {
		return UpdateRequest{}, err
	}

	req.Title = title
	req.Path = path
	req.Component = component
	req.Icon = icon
	req.Status = status
	req.Remark = remark
	return req, nil
}

func normalizeFields(menuType model.MenuType, title string, path string, component string, icon string, status model.MenuStatus, remark string) (string, string, string, string, model.MenuStatus, string, error) {
	if !ValidType(menuType) {
		return "", "", "", "", 0, "", errorsx.BadRequest("菜单类型不正确")
	}

	title = strings.TrimSpace(title)
	if title == "" {
		return "", "", "", "", 0, "", errorsx.BadRequest("菜单名称不能为空")
	}
	if len(title) > 64 {
		return "", "", "", "", 0, "", errorsx.BadRequest("菜单名称不能超过 64 个字符")
	}

	path = strings.TrimSpace(path)
	component = strings.TrimSpace(component)
	icon = strings.TrimSpace(icon)
	remark = strings.TrimSpace(remark)

	if len(path) > 255 {
		return "", "", "", "", 0, "", errorsx.BadRequest("路由路径不能超过 255 个字符")
	}
	if len(component) > 255 {
		return "", "", "", "", 0, "", errorsx.BadRequest("组件路径不能超过 255 个字符")
	}
	if len(icon) > 64 {
		return "", "", "", "", 0, "", errorsx.BadRequest("图标标识不能超过 64 个字符")
	}
	if len(remark) > 255 {
		return "", "", "", "", 0, "", errorsx.BadRequest("备注不能超过 255 个字符")
	}

	if status == 0 {
		status = model.MenuStatusEnabled
	}
	if !ValidStatus(status) {
		return "", "", "", "", 0, "", errorsx.BadRequest("菜单状态不正确")
	}

	if menuType == model.MenuTypeMenu && path == "" {
		return "", "", "", "", 0, "", errorsx.BadRequest("菜单节点需要填写路由路径")
	}

	return title, path, component, icon, status, remark, nil
}

func ValidType(menuType model.MenuType) bool {
	return menuType == model.MenuTypeDirectory ||
		menuType == model.MenuTypeMenu ||
		menuType == model.MenuTypeButton
}

func ValidStatus(status model.MenuStatus) bool {
	return status == model.MenuStatusEnabled || status == model.MenuStatusDisabled
}

func BuildResponse(item model.Menu) Response {
	return Response{
		ID:        item.ID,
		ParentID:  item.ParentID,
		Type:      item.Type,
		Code:      item.Code,
		Title:     item.Title,
		Path:      item.Path,
		Component: item.Component,
		Icon:      item.Icon,
		Sort:      item.Sort,
		Status:    item.Status,
		Remark:    item.Remark,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}
