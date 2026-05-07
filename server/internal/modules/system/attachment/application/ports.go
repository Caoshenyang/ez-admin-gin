package application

import (
	attachmentdomain "ez-admin-gin/server/internal/modules/system/attachment/domain"
	fileapp "ez-admin-gin/server/internal/modules/system/file/application"
	"ez-admin-gin/server/internal/platform/database"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

type AttachmentRepository interface {
	List(query attachmentdomain.ListQuery, page int, pageSize int) ([]attachmentdomain.View, int64, error)
	Create(tx *gorm.DB, item *attachmentdomain.Entity) error
	FindByID(tx *gorm.DB, id uint) (attachmentdomain.Entity, error)
	FindViewByID(id uint) (attachmentdomain.View, error)
	UpdateBase(tx *gorm.DB, item *attachmentdomain.Entity, req attachmentdomain.UpdateRequest) error
	UpdateStatus(tx *gorm.DB, item *attachmentdomain.Entity, status model.SystemAttachmentStatus) error
}

type AttachmentTransactor = database.Transactor
type FileAssetService = fileapp.AssetService
