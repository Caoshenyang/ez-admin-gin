package user

import (
	"ez-admin-gin/server/internal/platform/datascope"

	"gorm.io/gorm"
)

// applyDataScope 把用户资源的数据权限规则固定在一个地方，避免散落在 Handler 或 Service 里。
func applyDataScope(db *gorm.DB, actor datascope.Actor) *gorm.DB {
	return db.Scopes(datascope.UserQueryScope(db, actor, "department_id", "id"))
}
