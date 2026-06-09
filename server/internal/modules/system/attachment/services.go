package attachment

import (
	attachmentapp "ez-admin-gin/server/internal/modules/system/attachment/application"
	attachmentinfra "ez-admin-gin/server/internal/modules/system/attachment/infra"
	filemodule "ez-admin-gin/server/internal/modules/system/file"
	platformConfig "ez-admin-gin/server/internal/platform/config"
	platformDatabase "ez-admin-gin/server/internal/platform/database"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ServiceOptions struct {
	DB            *gorm.DB
	Upload        platformConfig.UploadConfig
	RuntimeConfig *platformConfig.RuntimeStore
	Log           *zap.Logger
}

// NewService 收拢 attachment 模块依赖装配，避免路由层重复拼接 repository / transactor / file asset service。
func NewService(opts ServiceOptions) *attachmentapp.Service {
	repo := attachmentinfra.NewRepository(opts.DB)
	transactor := platformDatabase.NewTransactor(opts.DB)
	fileService := filemodule.NewAssetService(filemodule.ServiceOptions{
		DB:            opts.DB,
		Upload:        opts.Upload,
		RuntimeConfig: opts.RuntimeConfig,
		Log:           opts.Log,
	})
	return attachmentapp.NewService(transactor, repo, fileService)
}
