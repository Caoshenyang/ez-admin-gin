package file

import (
	fileapp "ez-admin-gin/server/internal/modules/system/file/application"
	fileinfra "ez-admin-gin/server/internal/modules/system/file/infra"
	platformConfig "ez-admin-gin/server/internal/platform/config"
	platformDatabase "ez-admin-gin/server/internal/platform/database"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ServiceOptions 收拢 file 模块运行时依赖，避免调用方重复拼接 repository / storage / transactor。
type ServiceOptions struct {
	DB     *gorm.DB
	Upload platformConfig.UploadConfig
	Log    *zap.Logger
}

// NewService 构造 file 模块完整应用服务，供路由层直接使用。
func NewService(opts ServiceOptions) *fileapp.Service {
	repo := fileinfra.NewRepository(opts.DB)
	storage := fileinfra.NewLocalStorage(opts.Upload)
	transactor := platformDatabase.NewTransactor(opts.DB)
	return fileapp.NewService(transactor, repo, storage, opts.Upload, opts.Log)
}

// NewAssetService 暴露文件资产能力接口，供 attachment 等关联模块复用上传/回滚逻辑。
func NewAssetService(opts ServiceOptions) fileapp.AssetService {
	return NewService(opts)
}
