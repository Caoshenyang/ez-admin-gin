package setup

import (
	platformDatabase "ez-admin-gin/server/internal/platform/database"

	"gorm.io/gorm"
)

type ServiceOptions struct {
	DB *gorm.DB
}

// NewAppService 收拢 setup 模块依赖装配，让 routes 层只关心模块入口。
func NewAppService(opts ServiceOptions) *Service {
	repo := NewRepository(opts.DB)
	transactor := platformDatabase.NewTransactor(opts.DB)
	return NewService(transactor, repo)
}
