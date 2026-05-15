package infra

import (
	"errors"
	"time"

	notidomain "ez-admin-gin/server/internal/modules/system/notification/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

// Repository 封装通知表的数据访问操作。
type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// List 按用户、类型和已读状态分页查询通知列表。
func (r *Repository) List(userID uint, page int, pageSize int, notiType *model.NotificationType, isRead *bool) ([]notidomain.Entity, int64, error) {
	queryDB := r.db.Model(&notidomain.Entity{}).Where("user_id = ?", userID)

	if notiType != nil {
		queryDB = queryDB.Where("type = ?", *notiType)
	}
	if isRead != nil {
		queryDB = queryDB.Where("is_read = ?", *isRead)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []notidomain.Entity
	if err := queryDB.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// FindByID 在指定事务中按主键查找通知。
func (r *Repository) FindByID(db *gorm.DB, id uint64) (notidomain.Entity, error) {
	var item notidomain.Entity
	err := db.First(&item, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return notidomain.Entity{}, errorsx.NotFound("通知不存在")
		}
		return notidomain.Entity{}, err
	}
	return item, nil
}

// Create 在指定事务中插入一条新的通知记录。
func (r *Repository) Create(db *gorm.DB, item *notidomain.Entity) error {
	return db.Create(item).Error
}

// MarkRead 标记指定通知为已读。
func (r *Repository) MarkRead(db *gorm.DB, userID uint, ids []uint64) error {
	now := time.Now()
	return db.Model(&notidomain.Entity{}).
		Where("user_id = ? AND id IN ?", userID, ids).
		Updates(map[string]any{
			"is_read": true,
			"read_at": now,
		}).Error
}

// MarkAllRead 标记用户所有未读通知为已读。
func (r *Repository) MarkAllRead(db *gorm.DB, userID uint) error {
	now := time.Now()
	return db.Model(&notidomain.Entity{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]any{
			"is_read": true,
			"read_at": now,
		}).Error
}

// UnreadCount 获取用户未读通知数。
func (r *Repository) UnreadCount(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&notidomain.Entity{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}
