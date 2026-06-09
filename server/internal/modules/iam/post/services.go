package post

import (
	postapp "ez-admin-gin/server/internal/modules/iam/post/application"
	postinfra "ez-admin-gin/server/internal/modules/iam/post/infra"
	platformDatabase "ez-admin-gin/server/internal/platform/database"

	"gorm.io/gorm"
)

type ServiceOptions struct {
	DB *gorm.DB
}

// NewService 收拢岗位模块的依赖装配，减少 routes/api 对 infra 细节的感知。
func NewService(opts ServiceOptions) *postapp.Service {
	repo := postinfra.NewRepository(opts.DB)
	transactor := platformDatabase.NewTransactor(opts.DB)
	return postapp.NewService(transactor, repo)
}
