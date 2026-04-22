package permission

import (
	"fmt"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// Enforcer 包装 Casbin 权限判断能力。
type Enforcer struct {
	inner *casbin.Enforcer
}

// NewEnforcer 创建权限判断器，并从数据库加载策略。
func NewEnforcer(db *gorm.DB, modelPath string) (*Enforcer, error) {
	// 本项目统一使用 SQL 建表，不让 gorm-adapter 自动迁移表结构。
	gormadapter.TurnOffAutoMigrate(db)

	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("create casbin adapter: %w", err)
	}

	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		return nil, fmt.Errorf("create casbin enforcer: %w", err)
	}

	if err := enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("load casbin policy: %w", err)
	}

	return &Enforcer{
		inner: enforcer,
	}, nil
}

// Enforce 判断角色是否允许访问某个接口。
func (e *Enforcer) Enforce(sub string, obj string, act string) (bool, error) {
	allowed, err := e.inner.Enforce(sub, obj, act)
	if err != nil {
		return false, fmt.Errorf("enforce permission: %w", err)
	}

	return allowed, nil
}
