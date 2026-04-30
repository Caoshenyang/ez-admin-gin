package menu

import (
	"strings"
	"time"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
)

// CreateRequest 表示创建目录、菜单或按钮的请求体。
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

// UpdateRequest 表示编辑菜单基础信息的请求体。
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

// UpdateStatusRequest 表示单独更新菜单状态的请求体。
type UpdateStatusRequest struct {
	Status model.MenuStatus `json:"status"`
}

// Response 表示菜单树节点返回结构。
type Response struct {
	ID        uint               `json:"id"`
	ParentID  uint               `json:"parent_id"`
	Type      model.MenuType     `json:"type"`
	Code      string             `json:"code"`
	Title     string             `json:"title"`
	Path      string             `json:"path"`
	Component string             `json:"component"`
	Icon      string             `json:"icon"`
	Sort      int                `json:"sort"`
	Status    model.MenuStatus   `json:"status"`
	Remark    string             `json:"remark"`
	Children  []Response         `json:"children,omitempty"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// NormalizeCreateRequest 统一校验并收敛创建菜单参数。
func NormalizeCreateRequest(req CreateRequest) (CreateRequest, error) {
	code := strings.TrimSpace(req.Code)
	if code == "" {
		return CreateRequest{}, apperror.BadRequest("菜单编码不能为空")
	}
	if len(code) > 128 {
		return CreateRequest{}, apperror.BadRequest("菜单编码不能超过 128 个字符")
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

// NormalizeUpdateRequest 统一校验并收敛编辑菜单参数。
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
		return "", "", "", "", 0, "", apperror.BadRequest("菜单类型不正确")
	}

	title = strings.TrimSpace(title)
	if title == "" {
		return "", "", "", "", 0, "", apperror.BadRequest("菜单名称不能为空")
	}
	if len(title) > 64 {
		return "", "", "", "", 0, "", apperror.BadRequest("菜单名称不能超过 64 个字符")
	}

	path = strings.TrimSpace(path)
	component = strings.TrimSpace(component)
	icon = strings.TrimSpace(icon)
	remark = strings.TrimSpace(remark)

	if len(path) > 255 {
		return "", "", "", "", 0, "", apperror.BadRequest("路由路径不能超过 255 个字符")
	}
	if len(component) > 255 {
		return "", "", "", "", 0, "", apperror.BadRequest("组件路径不能超过 255 个字符")
	}
	if len(icon) > 64 {
		return "", "", "", "", 0, "", apperror.BadRequest("图标标识不能超过 64 个字符")
	}
	if len(remark) > 255 {
		return "", "", "", "", 0, "", apperror.BadRequest("备注不能超过 255 个字符")
	}

	if status == 0 {
		status = model.MenuStatusEnabled
	}
	if !ValidStatus(status) {
		return "", "", "", "", 0, "", apperror.BadRequest("菜单状态不正确")
	}

	if menuType == model.MenuTypeMenu && path == "" {
		return "", "", "", "", 0, "", apperror.BadRequest("菜单节点需要填写路由路径")
	}

	return title, path, component, icon, status, remark, nil
}

// ValidType 判断菜单节点类型是否合法。
func ValidType(menuType model.MenuType) bool {
	return menuType == model.MenuTypeDirectory ||
		menuType == model.MenuTypeMenu ||
		menuType == model.MenuTypeButton
}

// ValidStatus 判断菜单状态是否合法。
func ValidStatus(status model.MenuStatus) bool {
	return status == model.MenuStatusEnabled || status == model.MenuStatusDisabled
}

// BuildResponse 把菜单模型压成 API 返回结构。
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
