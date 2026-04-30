package authz

import (
	legacyPermission "ez-admin-gin/server/internal/permission"

	"gorm.io/gorm"
)

// Enforcer 复用现有 Casbin 封装，实现向 v2 authz 命名空间平滑迁移。
type Enforcer = legacyPermission.Enforcer

// NewEnforcer 创建接口权限判断器。
func NewEnforcer(db *gorm.DB, modelPath string) (*Enforcer, error) {
	return legacyPermission.NewEnforcer(db, modelPath)
}
