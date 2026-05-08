package user

import (
	userapp "ez-admin-gin/server/internal/modules/iam/user/application"
	userinfra "ez-admin-gin/server/internal/modules/iam/user/infra"
	platformDatabase "ez-admin-gin/server/internal/platform/database"

	"gorm.io/gorm"
)

type ServiceOptions struct {
	DB *gorm.DB
}

// NewService 收拢用户模块的依赖装配，让 routes 层只关心模块入口。
func NewService(opts ServiceOptions) *userapp.Service {
	repo := userinfra.NewRepository(opts.DB)
	transactor := platformDatabase.NewTransactor(opts.DB)
	return userapp.NewService(transactor, repo)
}
