package followup

import (
	"ez-admin-gin/server/internal/platform/datascope"

	"gorm.io/gorm"
)

// applyDataScope 把客户跟进资源的数据权限规则固定在一个地方。
func applyDataScope(db *gorm.DB, actor datascope.Actor) *gorm.DB {
	return db.Scopes(datascope.UserQueryScope(db, actor, "f.department_id", "f.owner_user_id"))
}
