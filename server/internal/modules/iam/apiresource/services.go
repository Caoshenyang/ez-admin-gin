package apiresource

import (
	apiresourceapp "ez-admin-gin/server/internal/modules/iam/apiresource/application"
	apiresourceinfra "ez-admin-gin/server/internal/modules/iam/apiresource/infra"

	"gorm.io/gorm"
)

type ServiceOptions struct {
	DB *gorm.DB
}

// NewService 收拢接口权限元数据模块的依赖装配。
func NewService(opts ServiceOptions) *apiresourceapp.Service {
	repo := apiresourceinfra.NewRepository(opts.DB)
	return apiresourceapp.NewService(repo)
}
