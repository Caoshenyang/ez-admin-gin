package application

import (
	noticedomain "ez-admin-gin/server/internal/modules/system/notice/domain"
	"ez-admin-gin/server/internal/platform/database"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

type NoticeTransactor = database.Transactor

type NoticeRepository interface {
	List(query noticedomain.ListQuery, page int, pageSize int, status *model.NoticeStatus) ([]noticedomain.Entity, int64, error)
	FindByID(db *gorm.DB, noticeID uint) (noticedomain.Entity, error)
	Create(db *gorm.DB, item *noticedomain.Entity) error
	UpdateBase(db *gorm.DB, item *noticedomain.Entity, req noticedomain.UpdateRequest) error
	UpdateStatus(db *gorm.DB, item *noticedomain.Entity, status model.NoticeStatus) error
	Delete(db *gorm.DB, item *noticedomain.Entity) error
}
