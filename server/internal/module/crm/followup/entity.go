package followup

import (
	"time"

	"ez-admin-gin/server/internal/model"
)

// Entity 复用 CRM 客户跟进模型，统一模块内部命名。
type Entity = model.CustomerFollowUp

// View 表示客户跟进列表与详情返回时的联合查询视图。
type View struct {
	ID             uint
	CustomerID     uint
	CustomerName   string
	DepartmentID   uint
	DepartmentName string
	OwnerUserID    uint
	OwnerUsername  string
	OwnerNickname  string
	FollowType     string
	Subject        string
	Content        string
	Result         string
	NextFollowAt   *time.Time
	Status         model.CustomerFollowUpStatus
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
