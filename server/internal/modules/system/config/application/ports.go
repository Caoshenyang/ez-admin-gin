package application

import (
	"context"
	"time"

	configdomain "ez-admin-gin/server/internal/modules/system/config/domain"
	"ez-admin-gin/server/internal/platform/database"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

type ConfigRepository interface {
	List(query configdomain.ListQuery, page int, pageSize int) ([]model.SystemConfig, int64, error)
	FindByID(db *gorm.DB, configID uint) (model.SystemConfig, error)
	FindEnabledByKey(key string) (model.SystemConfig, error)
	KeyExists(db *gorm.DB, key string) (bool, error)
	Create(db *gorm.DB, item *model.SystemConfig) error
	UpdateBase(db *gorm.DB, item *model.SystemConfig, req configdomain.UpdateRequest) error
	UpdateStatus(db *gorm.DB, item *model.SystemConfig, status model.SystemConfigStatus) error
}

type ConfigCache interface {
	Get(ctx context.Context, key string) (string, bool, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

type ConfigTransactor = database.Transactor
