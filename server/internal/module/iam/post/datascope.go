package post

import "gorm.io/gorm"

// applyDataScope 当前岗位资源不按部门范围裁剪，先保留显式落点，后续如果岗位归属到组织单元时再收紧规则。
func applyDataScope(db *gorm.DB) *gorm.DB {
	return db
}
