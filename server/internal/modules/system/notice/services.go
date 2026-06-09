package notice

import (
	noticeapp "ez-admin-gin/server/internal/modules/system/notice/application"
	noticeinfra "ez-admin-gin/server/internal/modules/system/notice/infra"
	platformDatabase "ez-admin-gin/server/internal/platform/database"

	"gorm.io/gorm"
)

type ServiceOptions struct {
	DB *gorm.DB
}

// NewService 收拢公告模块的依赖装配，避免 HTTP 路由层重复拼接 repository 和 transactor。
func NewService(opts ServiceOptions) *noticeapp.Service {
	repo := noticeinfra.NewRepository(opts.DB)
	transactor := platformDatabase.NewTransactor(opts.DB)
	return noticeapp.NewService(transactor, repo)
}
