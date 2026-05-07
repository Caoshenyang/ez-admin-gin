package infra

import (
	"ez-admin-gin/server/internal/platform/datascope"

	"gorm.io/gorm"
)

func applyDataScope(db *gorm.DB, actor datascope.Actor) *gorm.DB {
	return db.Scopes(datascope.UserQueryScope(db, actor, "department_id", "id"))
}
