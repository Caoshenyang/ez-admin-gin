package department

import (
	departmentapp "ez-admin-gin/server/internal/modules/iam/department/application"
	departmentinfra "ez-admin-gin/server/internal/modules/iam/department/infra"
	platformDatabase "ez-admin-gin/server/internal/platform/database"

	"gorm.io/gorm"
)

type ServiceOptions struct {
	DB *gorm.DB
}

// NewService 收拢部门模块的依赖装配，保持模块入口和 application service 的关系一致。
func NewService(opts ServiceOptions) *departmentapp.Service {
	repo := departmentinfra.NewRepository(opts.DB)
	transactor := platformDatabase.NewTransactor(opts.DB)
	return departmentapp.NewService(transactor, repo)
}
