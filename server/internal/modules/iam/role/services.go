package role

import (
	roleapp "ez-admin-gin/server/internal/modules/iam/role/application"
	roleinfra "ez-admin-gin/server/internal/modules/iam/role/infra"
	authzPlatform "ez-admin-gin/server/internal/platform/authz"
	platformDatabase "ez-admin-gin/server/internal/platform/database"

	"gorm.io/gorm"
)

type ServiceOptions struct {
	DB       *gorm.DB
	Enforcer *authzPlatform.Enforcer
}

// NewService 收拢角色模块的依赖装配，便于后续统一模块接入模板。
func NewService(opts ServiceOptions) *roleapp.Service {
	repo := roleinfra.NewRepository(opts.DB)
	transactor := platformDatabase.NewTransactor(opts.DB)
	return roleapp.NewService(transactor, repo, opts.Enforcer)
}
