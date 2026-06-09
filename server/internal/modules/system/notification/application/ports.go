package application

import (
	notidomain "ez-admin-gin/server/internal/modules/system/notification/domain"
	"ez-admin-gin/server/internal/platform/database"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

type NotificationTransactor = database.Transactor

type NotificationRepository interface {
	List(userID uint, page int, pageSize int, notiType *model.NotificationType, isRead *bool) ([]notidomain.Entity, int64, error)
	FindByID(db *gorm.DB, id uint64) (notidomain.Entity, error)
	Create(db *gorm.DB, item *notidomain.Entity) error
	MarkRead(db *gorm.DB, userID uint, ids []uint64) error
	MarkAllRead(db *gorm.DB, userID uint) error
	UnreadCount(userID uint) (int64, error)
}
