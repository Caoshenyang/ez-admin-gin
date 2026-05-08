package dict

import (
	dictapp "ez-admin-gin/server/internal/modules/system/dict/application"
	dictinfra "ez-admin-gin/server/internal/modules/system/dict/infra"
	platformDatabase "ez-admin-gin/server/internal/platform/database"

	"gorm.io/gorm"
)

type ServiceOptions struct {
	DB *gorm.DB
}

// NewService 收拢字典模块的依赖装配，保持模块入口和 application service 的衔接一致。
func NewService(opts ServiceOptions) *dictapp.Service {
	repo := dictinfra.NewRepository(opts.DB)
	transactor := platformDatabase.NewTransactor(opts.DB)
	return dictapp.NewService(transactor, repo)
}
