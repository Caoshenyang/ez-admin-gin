// Package domain 定义接口权限元数据的响应结构和业务规则。
package domain

import (
	"time"

	"ez-admin-gin/server/internal/platform/model"
)

// Response 定义接口权限元数据响应结构。
type Response struct {
	ID        uint            `json:"id"`
	Code      string          `json:"code"`
	Name      string          `json:"name"`
	Module    string          `json:"module"`
	Method    string          `json:"method"`
	Path      string          `json:"path"`
	Sort      int             `json:"sort"`
	Status    model.APIStatus `json:"status"`
	Remark    string          `json:"remark"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// Entity 是接口权限元数据模型的类型别名。
type Entity = model.API

const (
	PermissionList = "system:api:list"
)

// BuildResponse 将接口权限元数据模型转换为响应结构。
func BuildResponse(item model.API) Response {
	return Response{
		ID:        item.ID,
		Code:      item.Code,
		Name:      item.Name,
		Module:    item.Module,
		Method:    item.Method,
		Path:      item.Path,
		Sort:      item.Sort,
		Status:    item.Status,
		Remark:    item.Remark,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}
