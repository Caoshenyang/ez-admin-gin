package authz

import (
	"fmt"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

type Enforcer struct {
	inner *casbin.Enforcer
}

func NewEnforcer(db *gorm.DB, modelPath string) (*Enforcer, error) {
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

	return &Enforcer{inner: enforcer}, nil
}

func (e *Enforcer) Enforce(sub string, obj string, act string) (bool, error) {
	allowed, err := e.inner.Enforce(sub, obj, act)
	if err != nil {
		return false, fmt.Errorf("enforce permission: %w", err)
	}
	return allowed, nil
}
