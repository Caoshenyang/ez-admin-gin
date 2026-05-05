package customer

import (
	"time"

	"ez-admin-gin/server/internal/model"
)

// Entity 复用 CRM 客户模型，统一模块内部命名。
type Entity = model.Customer

// View 表示客户列表与详情返回时的联合查询视图。
type View struct {
	ID             uint
	Name           string
	ContactName    string
	Phone          string
	Level          string
	Source         string
	DepartmentID   uint
	DepartmentName string
	OwnerUserID    uint
	OwnerUsername  string
	OwnerNickname  string
	Status         model.CustomerStatus
	Remark         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
