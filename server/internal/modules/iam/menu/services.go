package menu

import (
	menuapp "ez-admin-gin/server/internal/modules/iam/menu/application"
	menuinfra "ez-admin-gin/server/internal/modules/iam/menu/infra"
	platformDatabase "ez-admin-gin/server/internal/platform/database"

	"gorm.io/gorm"
)

type ServiceOptions struct {
	DB *gorm.DB
}

// NewService 收拢菜单模块的依赖装配，减少 api/routes 对 infra 细节的感知。
func NewService(opts ServiceOptions) *menuapp.Service {
	repo := menuinfra.NewRepository(opts.DB)
	transactor := platformDatabase.NewTransactor(opts.DB)
	return menuapp.NewService(transactor, repo)
}
