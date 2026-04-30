package department

import (
	"ez-admin-gin/server/internal/platform/datascope"

	"gorm.io/gorm"
)

func applyDataScope(db *gorm.DB, actor datascope.Actor) *gorm.DB {
	return db.Scopes(datascope.DepartmentQueryScope(db, actor, "id"))
}
